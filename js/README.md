# infisical-sdk

A simple, zero-dependency Infisical SDK.

## Usage

### CJS

```js
const { InfisicalSDK } = require("infisical-sdk");

const client = new InfisicalSDK({
    siteUrl: "https://app.infisical.com" // optional
});

await client.login("<client-id>", "<client-secret>");
const secrets = await client.secrets("dev", "<project-id>");
```

### ESM

```js
import { InfisicalSDK } from "infisical-sdk";

const client = new InfisicalSDK({
    siteUrl: "https://app.infisical.com" // optional
});

await client.login("<client-id>", "<client-secret>");
const secrets = await client.secrets("dev", "<project-id>");
```