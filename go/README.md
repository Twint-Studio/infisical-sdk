# infisical-sdk

A simple, zero-dependency Infisical SDK for Go.

## Usage

```go
package main

import (
    "fmt"
    "log"

    "github.com/prince527/infisical-sdk"
)

func main() {
    sdk := infisical.New(infisical.Options{
        SiteURL: "https://app.infisical.com", // optional
    })

    if err := sdk.Login("<client-id>", "<client-secret>"); err != nil {
        log.Fatal(err)
    }

    secrets, err := sdk.Secrets("dev", "<project-id>")
    if err != nil {
        log.Fatal(err)
    }

    for _, s := range secrets.Secrets {
        fmt.Printf("%s = %s\n", s.SecretKey, s.SecretValue)
    }
}
```