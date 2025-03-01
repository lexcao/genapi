# GenAPI Landing Page Designs

This repository contains three distinct landing page designs for the GenAPI library, a declarative HTTP client generator for Go.

## Designs Overview

- **Design 1**: Clean, minimalist design with a professional and modern look
- **Design 2**: Dark theme with a more technical aesthetic 
- **Design 3**: Vibrant, colorful design with gradients and modern UI effects

Each design is built using vanilla HTML, CSS, and JavaScript without any frameworks to ensure GitHub Pages compatibility.

## Local Development

To run the designs locally, you can use the included development server:

```bash
# Navigate to the shared directory
cd shared

# Run the server (defaults to port 8080)
go run server.go

# Run with custom port
go run server.go -port 3000

# Run from a different directory
go run server.go -dir "../" -port 8080
```

This will start a local server, and you can access the designs at:

- Design 1: http://localhost:8080/design1/
- Design 2: http://localhost:8080/design2/
- Design 3: http://localhost:8080/design3/

## Directory Structure

```
website/
├── shared/         # Shared resources and server
│   └── server.go   # Simple Go HTTP server for local development
├── design1/        # Clean, minimalist design
│   ├── index.html
│   ├── styles.css
│   └── script.js
├── design2/        # Dark theme, technical design
│   ├── index.html
│   ├── styles.css
│   └── script.js
├── design3/        # Colorful, gradient-based design
│   ├── index.html
│   ├── styles.css
│   └── script.js
├── api/            # Sample API code files
│   ├── todo.go     # Original API definition
│   └── todo.gen.go # Generated client code
└── README.md       # This file
```

## Features

Each design includes:

- Responsive layout for all device sizes
- Interactive code examples with syntax highlighting
- Tabbed interface to switch between code examples
- Mobile-friendly navigation
- Smooth scrolling for anchor links
- Hover effects and animations for improved user engagement

## Customization

You can easily customize these designs by:

1. Modifying the CSS in the respective `styles.css` files
2. Updating the content in the `index.html` files
3. Extending the JavaScript functionality in the `script.js` files

## Deployment

To deploy to GitHub Pages:

1. Push the content to a GitHub repository
2. Go to repository Settings → Pages
3. Select the branch and folder containing the website
4. GitHub will provide a URL to your published site

## Credits

- Code syntax highlighting: [PrismJS](https://prismjs.com/)
- Icons: SVG from various sources 