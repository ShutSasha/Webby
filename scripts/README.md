# Secrets Management Setup

A script was added in this branch to **automatically inject secrets into configuration files**.  
Secrets are stored using **Infisical**.

The script fetches secrets from Infisical and replaces placeholders in `.temp` configuration files.

---

# 1. Install Infisical

Install the CLI on your machine.

### Windows
```powershell
winget install infisical
```

After installation, log in:
```powershell
infisical login
```

Use the EU server and the project mailsystem account.

---

# 2. Generate Configuration Files

To generate configuration files with injected secrets:

- Go to the scripts directory

- Run the script

```shell
python generate_config.py
```

The script will:

- fetch secrets from Infisical

- replace placeholders in .temp files

- generate final configuration files

Example:

- ```appsettings.Development.json.temp → appsettings.Development.json```
- ```config.yaml.temp → config.yaml```
- ```.env.temp → .env```

Inside the script you must add paths to your .temp files in the configs list.

Example:
```py
configs = [
    "..\\server\\Webby.AuthService\\appsettings.Development.json.temp",
    "..\\server\\Webby.UserService\\appsettings.Development.json.temp",
    "..\\server\\Webby.ApiGetaway\\appsettings.Development.json.temp",
    "..\\web\\.env.temp"
]
```
---

# 3. Adding Secrets

Secrets are added through the Infisical web dashboard.

Steps:

- Log in to the Infisical account

- Open the project

- Add a new secret

---

# 4. Using Secrets in Config Files

Create configuration template files with the .temp extension.

Examples:

- ```.env.temp```
- ```config.yaml.temp```
- ```appsettings.Development.json.temp```

Secrets are referenced using the following format:

```{{SecretName}}```

Example .env.temp:
```
AUTH_SECRET={{WebAuthSecret}}
NEXT_PUBLIC_API_URL=http://localhost:5000/api
AUTH_GOOGLE_ID={{WebAuthGoogleId}}
AUTH_GOOGLE_SECRET={{WebGoogleAuthSecret}}
```

When the script runs, placeholders will be automatically replaced with the actual values stored in Infisical.

---

# Summary

1. Install Infisical

2. Login to the project account

3. Add .temp config files

4. Add secret names using {{SecretName}}

5. Run the script to generate real config files