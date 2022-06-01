package h2c_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

func TestBN254G1Nu(t *testing.T) {
	jsonFile, err := os.Open("testdata/suites/BN254G1_XMD-SHA-256_SVDW_NU_.json")
	if err != nil {
		t.Fatal(err)
	}
	defer jsonFile.Close()
	var bytes []byte
	if bytes, err = ioutil.ReadAll(jsonFile); err != nil {
		t.Fatal(err)
	}
	var v vectorSuite
	if err = json.Unmarshal(bytes, &v); err != nil {
		t.Fatal(err)
	}
	v.test(t)
}

func TestBLS12_381G2Nu(t *testing.T) {
	jsonFile, err := os.Open("testdata/suites/BLS12381G2_XMD-SHA-256_SSWU_NU_.json")
	if err != nil {
		t.Fatal(err)
	}
	defer jsonFile.Close()
	var bytes []byte
	if bytes, err = ioutil.ReadAll(jsonFile); err != nil {
		t.Fatal(err)
	}
	var v vectorSuite
	if err = json.Unmarshal(bytes, &v); err != nil {
		t.Fatal(err)
	}
	v.test(t)
}
