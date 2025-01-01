# Your awesome project

Welcome to the yet another awesome project generated with the Gowebly CLI! 🎉

This README file contains all the necessary information about:

- [Project overview](#project-overview)
- [Folders structure](#folders-structure)
- [Starting your project](#starting-your-project)
- [Developing your project](#developing-your-project)
- [Deploying your project](#deploying-your-project)

And some words [about the Gowebly CLI](#about-the-gowebly-cli).

## Project overview

Backend:

- Module name in the go.mod file: `fmj`
- Go web framework/router: `Gin`
- Server port: `7000`

Frontend:

- Package name in the package.json file: `fundmyjollof`
- Reactivity library: `htmx with hyperscript`
- CSS framework: `Tailwind CSS with Preline UI components`

Tools:

- Air tool to live-reloading: ✓
- Bun as a frontend runtime: ✓
- Templ to generate HTML: ✕
- Config for golangci-lint: ✓

## Folders structure

```console
.
├── assets
│   ├── scripts.js
│   └── styles.scss
├── static
│   ├── images
│   │   └── gowebly.svg
│   ├── apple-touch-icon.png
│   ├── favicon.ico
│   ├── favicon.png
│   ├── favicon.svg
│   ├── manifest-desktop-screenshot.jpeg
│   ├── manifest-mobile-screenshot.jpeg
│   ├── manifest-touch-icon.svg
│   └── manifest.webmanifest
├── templates
│   ├── pages
│   │   └── index.html
│   └── main.html
├── .gitignore
├── .dockerignore
├── .prettierignore
├── .air.toml
├── golangci.yml
├── Dockerfile
├── docker-compose.yml
├── prettier.config.js
├── package.json
├── go.mod
├── go.sum
├── handlers.go
├── server.go
├── main.go
└── README.md
```

## Starting your project

> ❗️ Please make sure that you have installed the executable files for all the necessary tools before starting your project. Exactly:
> - `Air`: [https://github.com/air-verse/air](https://github.com/air-verse/air)
> - `Bun`: [https://github.com/oven-sh/bun](https://github.com/oven-sh/bun)

> - `golangci-lint`: [https://github.com/golangci/golangci-lint](https://github.com/golangci/golangci-lint)

To start your project, run the **Gowebly** CLI command in your terminal:

```console
gowebly run
```

## Developing your project

The backend part is located in the `*.go` files in your project folder.

The `./templates` folder contains HTML templates that you can use in your frontend part. Also, the `./assets` folder contains the `styles.scss` (main styles) and `scripts.js` (main scripts) files.

The `./static` folder contains all the static files: icons, images, PWA (Progressive Web App) manifest and other builded/minified assets.

## Deploying your project

All deploy settings are located in the `Dockerfile` and `docker-compose.yml` files in your project folder.

To deploy your project to a remote server, follow these steps:

1. Go to your hosting/cloud provider and create a new VDS/VPS.
2. Update all OS packages on the server and install Docker, Docker Compose and Git packages.
3. Use `git clone` command to clone the repository with your project to the server and navigate to its folder.
4. Run the `docker-compose up` command to start your project on your server.



## About the Gowebly CLI

The [**Gowebly**](https://github.com/gowebly/gowebly) CLI is a next-generation CLI tool that makes it easy to create amazing web applications with **Go** on the backend, using **htmx**, **hyperscript** or **Alpine.js**, and the most popular **CSS frameworks** on the frontend.

It's highly recommended to start exploring the Gowebly CLI with short articles "[**What is Gowebly CLI?**](https://gowebly.org/getting-started)" and "[**How does it work?**](https://gowebly.org/getting-started/how-does-it-work)" to understand the basic principle and the main components built into the **Gowebly** CLI.

<a href="https://gowebly.org/" target="_blank"><img height="112px" alt="another awesome project built with the Gowebly CLI" src="https://raw.githubusercontent.com/gowebly/.github/main/images/gowebly-new-project-banner.svg"/></a>
