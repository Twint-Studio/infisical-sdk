class InfisicalSDK {
    constructor(options = {}) {
        this.siteUrl = (options.siteUrl || "https://app.infisical.com").replace(/\/+$/, "");
        this.accessToken = null;
    }

    async #request(path, { params, ...init } = {}) {
        const url = new URL(path, this.siteUrl);
        if (params) {
            for (const [key, value] of Object.entries(params)) {
                url.searchParams.set(key, value);
            }
        }

        const res = await fetch(url, {
            ...init,
            headers: {
                ...(this.accessToken && { Authorization: `Bearer ${this.accessToken}` }),
                ...init.headers
            }
        });

        const data = await res.json().catch(() => ({}));
        if (!res.ok) throw new Error(data.message || data.error || `HTTP ${res.status}`);

        return data;
    }

    async login(id, secret) {
        const data = await this.#request("/api/v1/auth/universal-auth/login", {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                clientId: id,
                clientSecret: secret
            })
        });

        this.accessToken = data.accessToken;
    }

    async secrets(environment, project) {
        if (!this.accessToken) throw new Error("Not authenticated. Call login() first.");

        return this.#request("/api/v3/secrets/raw", {
            params: {
                environment,
                workspaceId: project
            }
        });
    }
}

export {
    InfisicalSDK
};