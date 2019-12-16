# Pingboard (work in progress)
## What is this
This small project is about to visualize basic connectivity/health status of a small set of network/computer equipments.
The idea is not about to create a complex monitoring system (since those are already exist, right? :)), but create something simple, clean and flexible for ad-hoc checking.

## How it will work
The project is going to be written in Go. I don't like frontend programming, so the aim is to program as less client-side code or any HTML as possible. Instead, the idea is to handcraft an SVG graphics file for your use-case (map with the equipments, cables, etc.), which will be re-colored according to the last status by this project and served via HTTP. Since SVG is essentially XML and all drawn things must have unique IDs, it is easy to change properties of objects from SW.
SVG moreover can handle alt-texts and several other useful things for a good user experience.
There are going to be multiple status queries, like ICMP ping, HTTP GET, custom script execution, etc..

The idea is to create go code which is as idiomatic as possible and according to the best practices. This means prefering communication over memory share (i.e. avoid mutexes), etc.

## Configuration
The (YAML based) configuration will refer to the actual SVG and further references for the objects in it, along with the associated healthcheck. The user will see the on-demand rendered SVG with the latest statuses. See `examples` folder for sample.