![svg-to-dxf](logo.png)

> Convert .svg files to .dxf

Available at https://svg-to-dxf.cluster.fun/

## Features

Runs a webserver that takes in an `?url=` query string and fetches the SVG from that URL or an svg file uploaded and then returns it as a .dxf

## Usage

```sh
docker run -it --rm -p 8080:8080 docker.cluster.fun/averagemarcus/svg-to-dxf
```

## Building from source

With Docker:

```sh
make docker-build
```

## Resources

* [inkscape](https://inkscape.org/)

## Contributing

If you find a bug or have an idea for a new feature please raise an issue to discuss it.

Pull requests are welcomed but please try and follow similar code style as the rest of the project and ensure all tests and code checkers are passing.

Thank you 💛

## License

See [LICENSE](LICENSE)
