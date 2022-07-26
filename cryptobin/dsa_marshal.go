package cryptobin

import (
    "fmt"
    "math/big"
    "crypto/dsa"
    "encoding/asn1"
)

// 序列号
var dsaPrivKeyVersion = 1

// 私钥
type dsaPrivateKey struct {
    Version int
    P       *big.Int
    Q       *big.Int
    G       *big.Int
    Y       *big.Int
    X       *big.Int
}

// 公钥
type dsaPublicKey struct {
    P *big.Int
    Q *big.Int
    G *big.Int
    Y *big.Int
}

// 包装公钥
func (this DSA) MarshalPublicKey(key *dsa.PublicKey) ([]byte, error) {
    publicKey := dsaPublicKey{
        P: key.Parameters.P,
        Q: key.Parameters.Q,
        G: key.Parameters.G,
        Y: key.Y,
    }

    return asn1.Marshal(publicKey)
}

// 解析公钥
func (this DSA) ParsePublicKey(derBytes []byte) (*dsa.PublicKey, error) {
    var key dsaPublicKey
    rest, err := asn1.Unmarshal(derBytes, &key)
    if err != nil {
        return nil, err
    }

    if len(rest) > 0 {
        return nil, asn1.SyntaxError{Msg: "trailing data"}
    }

    publicKey := &dsa.PublicKey{
        Parameters: dsa.Parameters{
            P: key.P,
            Q: key.Q,
            G: key.G,
        },
        Y: key.Y,
    }

    return publicKey, nil
}

// ====================

// 包装私钥
func (this DSA) MarshalPrivateKey(key *dsa.PrivateKey) ([]byte, error) {
    // 版本号
    version := dsaPrivKeyVersion

    // 公钥
    publicKey := key.PublicKey

    // 构造私钥信息
    privateKey := dsaPrivateKey{
        Version: version,
        P:       publicKey.Parameters.P,
        Q:       publicKey.Parameters.Q,
        G:       publicKey.Parameters.G,
        Y:       publicKey.Y,
        X:       key.X,
    }

    return asn1.Marshal(privateKey)
}

// 解析私钥
func (this DSA) ParsePrivateKey(derBytes []byte) (*dsa.PrivateKey, error) {
    var key dsaPrivateKey
    rest, err := asn1.Unmarshal(derBytes, &key)
    if err != nil {
        return nil, err
    }

    if len(rest) > 0 {
        return nil, asn1.SyntaxError{Msg: "trailing data"}
    }

    if key.Version != dsaPrivKeyVersion {
        return nil, fmt.Errorf("DSA: unknown DSA private key version %d", key.Version)
    }

    privateKey := &dsa.PrivateKey{
        PublicKey: dsa.PublicKey{
            Parameters: dsa.Parameters{
                P: key.P,
                Q: key.Q,
                G: key.G,
            },
            Y: key.Y,
        },
        X: key.X,
    }

    return privateKey, nil
}