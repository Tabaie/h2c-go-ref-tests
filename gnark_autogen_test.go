package h2c_test

import (
	"fmt"
	C "github.com/armfazh/tozan-ecc/curve"
	"github.com/armfazh/tozan-ecc/field"
	h2c "h2c-go-ref"
	"os"
	"os/exec"
	"strconv"
	"testing"
)

var msgs = []string{
	"",
	"abc", "abcdef0123456789", "q128_qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq", "a512_aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
}

func writeStr(f *os.File, s string) {
	n, err := f.WriteString(s)
	if err != nil {
		panic(err)
	}
	if n != len(s) {
		panic("not all of string written to file")
	}
}

func writeName(f *os.File, name string) {
	if name != "" {
		writeStr(f, name)
		writeStr(f, ": ")
	}
}

func writeElt(f *os.File, e field.Elt, name string) {
	writeName(f, name)
	writeStr(f, fmt.Sprintf("\"%v\"", e))
}

func writePoint(f *os.File, p C.Point, name string) {
	writeName(f, name)
	writeStr(f, "point{")
	writeElt(f, p.X(), "")
	writeStr(f, ", ")
	writeElt(f, p.Y(), "")
	writeStr(f, "}")
}

func generateGnarkSubVector(f *os.File, suiteId string, dst string) {
	hashToCurve, err := h2c.SuiteID(suiteId).Get([]byte(dst))
	if err != nil {
		panic(err)
	}

	for _, msg := range msgs {
		msgBytes := []byte(msg)

		u := hashToCurve.U(msgBytes)
		q := hashToCurve.Q(u)
		p := hashToCurve.Hash(msgBytes)

		writeStr(f,"{\nmsg: \"")
		writeStr(f,msg)
		writeStr(f,"\", ")

		writePoint(f,p, "P")
		writeStr(f,",\n")

		for i := 0; i < len(q); i++ {
			name := "Q"
			if len(q) != 1 {
				name += strconv.Itoa(i)
			}

			writePoint(f, q[i], name)
			writeStr(f, ",\n")
		}

		for i := 0; i < len(u); i++ {
			name := "u"
			if len(u) != 1 {
				name += strconv.Itoa(i)
			}

			writeElt(f, u[i], name)
			writeStr(f, ", ")
		}

		writeStr(f, "\n},")
	}
}

func writeTestVector(f *os.File, suiteIdTemplate string, dstTemplate string, randomOracle bool, curveName string) {
	var hashOrEncode string
	var randomOracleOrNonUniform string
	if randomOracle {
		hashOrEncode, randomOracleOrNonUniform =
			 "hash", "RO_"
	} else {
		hashOrEncode, randomOracleOrNonUniform =
			"encode", "NU_"
	}

	dst := fmt.Sprintf(dstTemplate, curveName) + randomOracleOrNonUniform
	suiteId :=  fmt.Sprintf(suiteIdTemplate, curveName) + randomOracleOrNonUniform

	writeStr(f, hashOrEncode)
	writeStr(f, "To")
	writeStr(f, curveName)
	writeStr(f, "Vector")
	writeStr(f," = ")
	writeStr(f, hashOrEncode)
	writeStr(f, "TestVector{\ndst: []byte(\"")
	writeStr(f, dst)
	writeStr(f,"\"), \n cases: []")
	writeStr(f, hashOrEncode)
	writeStr(f,"TestCase{\n")

	generateGnarkSubVector(f, suiteId, dst )

	writeStr(f,"\n}}\n")
}

func generateGnarkVector(suiteIdTemplate string, dstTemplate string, outFileName string) {

	f, err := os.Create(outFileName)
	if err != nil {
		panic(err)
	}


	writeStr(f,"package bn254\nfunc init() {")

	writeTestVector(f, suiteIdTemplate, dstTemplate, false, "G1")
	writeTestVector(f, suiteIdTemplate, dstTemplate, true, "G1")


	writeStr(f, "}\n")
	f.Close()
	cmd := exec.Command("gofmt", "-w", outFileName)
	cmd.Run()
}

func TestGenerateGnarkBn254Vector(t *testing.T) {
	generateGnarkVector(
		"BN254%s_XMD:SHA-256_SVDW_",
		"QUUX-V01-CS02-with-BN254%s_XMD:SHA-256_SVDW_",
		"../gnark-crypto/ecc/bn254/hash_vectors_test.go")
}