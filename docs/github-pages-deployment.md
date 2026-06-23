# Deploying the SDK Documentation to GitHub Pages

This repository includes a static landing page at:

```text
docs/index.html
```

GitHub Pages can serve this directly.

## Option A: Deploy from `/docs` folder

1. Push the repository to GitHub.
2. Open the repository on GitHub.
3. Go to **Settings → Pages**.
4. Under **Build and deployment**, choose:
   - Source: **Deploy from a branch**
   - Branch: `main` or `master`
   - Folder: `/docs`
5. Save.

GitHub will publish the page at:

```text
https://<your-github-username>.github.io/go-vodacom-sdk/
```

## Option B: Custom domain

1. Add your domain in **Settings → Pages → Custom domain**.
2. Add the required DNS record:
   - `CNAME` for a subdomain, e.g. `sdk.example.com`
   - `A` records for apex domain, per GitHub Pages docs
3. Add a `CNAME` file inside `docs/` if needed:

```text
sdk.example.com
```

## Local preview

From repository root:

```bash
python3 -m http.server 8080 -d docs
```

Then open:

```text
http://localhost:8080
```

## Notes

- The page is static and works without a backend.
- The file includes optional HTMX attributes for progressive enhancement.
- The visual design uses Vodacom-inspired red and white styling.
