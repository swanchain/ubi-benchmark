* [LoginByApikey](#LoginByApikey)
* [Register](#Register)
* [LoginByPublicKeySignature](#LoginByPublicKeySignature)
* [CheckLogin](#CheckLogin)
* [GenerateApikey](#GenerateApikey)
* [DeleteApikey](#DeleteApikey)
* [GetApikeys](#GetApikeys)
* [RegisterEmail](#RegisterEmail)
* [GetWallet](#GetWallet)
* [SetPopupTime](#SetPopupTime)
* [DeleteEmail](#DeleteEmail)


## LoginByApikey

Definition:

```shell
func LoginByApikey(apikey, accessToken, network string) (*McsClient, error)
```

Outputs:

```shell
*McsClient  # includes jwt token and other information for use when call the other apis
error       # error or nil
```

## Register

Definition:

```shell
func Register(publicKeyAddress, network string) (*string, error)
```

Outputs:

```shell
*string  # nonce used for login by public key
error    # error or nil
```


## LoginByPublicKeySignature

Definition:

```shell
func LoginByPublicKeySignature(nonce, publicKeyAddress, signature, network string) (*McsClient, error)
```

Outputs:

```shell
*McsClient  # includes jwt token and other information for use when call the other apis
error       # error or nil
```

## CheckLogin

Definition:

```shell
func (mcsClient *McsClient) CheckLogin() (*string, *string, error)
```

Outputs:

```shell
*string  # network name
*string  # wallet address
error    # error or nil
```

## GenerateApikey

Definition:

```shell
func (mcsClient *McsClient) GenerateApikey(validDys int) (*string, *string, error)
```

Outputs:

```shell
*string  # apikey
*string  # access token
error    # error or nil
```

## DeleteApikey

Definition:

```shell
func (mcsClient *McsClient) DeleteApikey(apikey string) error
```

Outputs:

```shell
error    # error or nil
```

## GetApikeys

Definition:

```shell
func (mcsClient *McsClient) GetApikeys() ([]*Apikey, error)
```

Outputs:

```shell
[]*Apikey # apikey lists that belong the current user
error     # error or nil
```

## RegisterEmail

Definition:

```shell
func (mcsClient *McsClient) RegisterEmail(email string) (*string, error)
```

Outputs:

```shell
*string   # illustrate message
error     # error or nil
```

## GetWallet

Definition:

```shell
func (mcsClient *McsClient) GetWallet() (*Wallet, error)
```

Outputs:

```shell
*Wallet   # wallet information
error     # error or nil
```

## SetPopupTime

Definition:

```shell
func (mcsClient *McsClient) SetPopupTime() error
```

Outputs:

```shell
error     # error or nil
```

## DeleteEmail

Definition:

```shell
func (mcsClient *McsClient) DeleteEmail() error
```

Outputs:

```shell
error     # error or nil
```

