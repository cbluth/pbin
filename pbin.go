package pbin

import (
	"bytes"
	"compress/flate"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/gearnode/base58"
	"golang.org/x/crypto/pbkdf2"
)

var (
	OpenDiscussion   bool
	BurnAfterReading bool
)

const (
	APIVersion          int    = 2
	Iterations          int    = 100000 // kdf iterations
	KDFSecretSize       int    = 32     // bytes
	AESKeySize          int    = 32     // bytes
	NonceSize           int    = 12     // bytes
	SaltSize            int    = 8      // bytes
	TagSize             int    = 128    // bits
	EncryptionAlgorithm string = "aes"
	EncryptionMode      string = "gcm"
	DataCompression     string = "zlib"
	Format              string = "syntaxhighlighting"
	Expiry              string = "1week"
	// OpenDiscussion      bool   = false
	// BurnAfterReading    bool   = false
)

type (
	Paste struct {
		Version             int // 2
		ClearTextData       []byte
		ClearJSONData       []byte
		CipherJSONData      []byte
		RequestBodyJSONData []byte
		KDFSecret           [KDFSecretSize]byte
		AESKey              [AESKeySize]byte
		Salt                [SaltSize]byte
		Nonce               [NonceSize]byte // IV
		Expire              string
		OpenDiscussion      bool
		BurnAfterReading    bool
		DisplayFormat       string
	}
)

func CraftPaste(bytes []byte) (*Paste, error) {
	p := &Paste{
		Version:       APIVersion,
		DisplayFormat: Format,
		ClearTextData: bytes,
	}
	copy(p.Salt[:], randomBytes(SaltSize))
	copy(p.Nonce[:], randomBytes(NonceSize)) // IV
	copy(p.KDFSecret[:], randomBytes(KDFSecretSize))
	p.Expire = Expiry
	p.DisplayFormat = Format
	p.OpenDiscussion = OpenDiscussion
	p.BurnAfterReading = BurnAfterReading
	err := p.encrypt()
	if err != nil {
		return nil, err
	}
	req := map[string]interface{}{}
	req["v"] = APIVersion
	req["adata"] = p.makeAData()
	req["meta"] = map[string]interface{}{}
	req["meta"].(map[string]interface{})["expire"] = Expiry
	req["ct"] = base64.RawStdEncoding.EncodeToString(p.CipherJSONData)
	p.RequestBodyJSONData, err = json.Marshal(&req)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Paste) Send() (*url.URL, map[string]interface{}, error) {
	host := findFastest()
	req, err := http.NewRequest(http.MethodPost, host.URL.String(), bytes.NewBuffer(p.RequestBodyJSONData))
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("X-Requested-With", "JSONHttpRequest")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, nil, errors.New("error from server: " + host.URL.String())
	}
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, nil, err
	}
	resm := map[string]interface{}{}
	err = json.Unmarshal(resBody, &resm)
	if err != nil {
		return nil, nil, err
	}
	if resm["status"].(float64) != 0 {
		return nil, nil, errors.New("error from server: " + resm["message"].(string))
	}
	purl, err := url.Parse(host.URL.String() + "?" + resm["id"].(string) + "#" + base58.Encode(p.KDFSecret[:]))
	if err != nil {
		return nil, nil, err
	}
	return purl, resm, nil
}

func randomBytes(n int) []byte {
	k := make([]byte, n)
	_, err := rand.Read(k[:n])
	if err != nil {
		panic(err)
	}
	return k
}

func (p *Paste) encrypt() error {
	err := (error)(nil)
	p.ClearJSONData, err = json.Marshal(
		&map[string]interface{}{
			"paste": string(p.ClearTextData),
		},
	)
	if err != nil {
		return err
	}
	copy(p.AESKey[:], makeAESKey(p.KDFSecret[:], p.Salt[:]))
	c, err := aes.NewCipher(p.AESKey[:])
	if err != nil {
		return err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return err
	}
	adata, err := json.Marshal(p.makeAData())
	if err != nil {
		return err
	}
	b := bytes.Buffer{}
	w, err := flate.NewWriter(&b, flate.BestCompression)
	if err != nil {
		return err
	}
	_, err = w.Write(p.ClearJSONData)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	p.CipherJSONData = gcm.Seal(nil, p.Nonce[:], b.Bytes(), adata)
	return nil
}

func (p *Paste) makeAData() []interface{} {
	openDiscussion := int(0)
	burnAfterRead := int(0)
	if p.OpenDiscussion {
		openDiscussion = 1
	}
	if p.BurnAfterReading {
		burnAfterRead = 1
	}
	return []interface{}{
		[]interface{}{
			base64.RawStdEncoding.EncodeToString(p.Nonce[:]), // IV
			base64.RawStdEncoding.EncodeToString(p.Salt[:]),  // salt
			Iterations,
			256,
			TagSize,
			EncryptionAlgorithm,
			EncryptionMode,
			DataCompression,
		},
		Format,
		openDiscussion,
		burnAfterRead,
	}
}

func makeAESKey(secret []byte, salt []byte) []byte {
	return pbkdf2.Key(
		secret,
		salt,
		Iterations,
		AESKeySize,
		sha256.New,
	)
}

func GetPaste(ur *url.URL) ([]byte, error) {
	pID := ur.RawQuery
	b58Pass := ur.Fragment
	hostURL := strings.Split(ur.String(), "?")[0]
	pasteDataURL := hostURL + "?pasteid=" + pID
	req, err := http.NewRequest(http.MethodGet, pasteDataURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Requested-With", "JSONHttpRequest")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	m := map[string]interface{}{}
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}
	p := &Paste{}
	if v, ok := m["ct"]; !ok {
		return nil, errors.New("missing ct")
	} else {
		p.CipherJSONData, err = base64.RawStdEncoding.DecodeString(v.(string))
		if err != nil {
			return nil, err
		}
	}
	if v, ok := m["adata"]; !ok {
		return nil, errors.New("missing adata")
	} else {
		nonceData, err := base64.RawStdEncoding.DecodeString(((v.([]interface{})[0]).([]interface{})[0]).(string)) // wtf
		if err != nil {
			return nil, err
		}
		copy(p.Nonce[:], nonceData)
		saltData, err := base64.RawStdEncoding.DecodeString(((v.([]interface{})[0]).([]interface{})[1]).(string)) // wtf
		if err != nil {
			return nil, err
		}
		copy(p.Salt[:], saltData)
	}
	secret, err := base58.Decode(b58Pass)
	if err != nil {
		return nil, err
	}
	aesKey := makeAESKey(secret, p.Salt[:])
	c, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}
	adata, err := json.Marshal(p.makeAData())
	if err != nil {
		return nil, err
	}
	flated, err := gcm.Open(nil, p.Nonce[:], p.CipherJSONData, adata)
	if err != nil {
		return nil, err
	}
	fr := flate.NewReader(bytes.NewBuffer(flated))
	defer fr.Close()
	unflated, err := ioutil.ReadAll(fr)
	if err != nil {
		return nil, err
	}
	pd := map[string]interface{}{}
	err = json.Unmarshal(unflated, &pd)
	if err != nil {
		return nil, err
	}
	if v, ok := pd["paste"]; ok {
		return []byte(v.(string)), nil
	}
	return nil, errors.New("missing paste data")
}
