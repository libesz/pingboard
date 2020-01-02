# Pingboard (work in progress)
## What is this
This small project is about to visualize basic connectivity/health status of a small set of network/computer equipments.
The idea is not about to create a complex monitoring system (since those are already exist, right? :)), but create something simple, clean and flexible for ad-hoc checking.

## How it works
The project is going to be written in Go. I don't like frontend programming (sorry), so the aim is to program as less client-side code or any HTML as possible. Instead, the idea is to handcraft an SVG graphics file for your use-case (map with the equipments, cables, etc.), which will be re-colored according to the last status by this project and served via HTTP/WebSocket. Since SVG is essentially XML and all drawn things can have unique IDs, it is easy to change properties of objects from SW and update the browser.
SVG moreover can handle alt-texts and several other useful things for a good user experience.
There are going to be multiple status queries, like ICMP ping, HTTP GET, custom script execution, etc. Currently Ping is implemented.

The idea is to create go code which is as idiomatic as possible and according to the best practices. This means:
* No platform specific solution
* Single binary (static web content is built-in)
* Prefering communication over shared state (i.e. avoid mutexes)

## Build
```
go build .
```

## Build after change in the static HTML content
```
go generate .
go build .
```

## Configuration
The (YAML based) configuration will refer to the actual SVG and further references for the objects in it, along with the associated healthcheck. The user is expected to provide the SVG document. The objects (shapes, paths, etc.) shall have `id` attributes, which will be searched based on the configuration. If you end up with an SVG without `id` properties, you shall fix it via an editor, like Inkscape.
The user will see the on-demand rendered SVG with the latest statuses. See `examples` folder for sample. After you have a valid config, run the service via: `./pingboard config.yaml`. Currently it listens on localhost:2003 and does not tolerate reverse proxy, but it will later...

## TODO
* Implement more checks
* Dockerize
* Tolerate reverse proxy / TLS
* Prometheus exporter
