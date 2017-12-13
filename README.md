# Running the Demo

Set up a cluster of `amd64`, `arm32v7`, and/or `arm64v8` machines running Docker.  Configure Docker Swarm Mode in the configuration of your choosing (keeping in mind that manager nodes are intentionally excluded from the demo in the provided `stack.yml`, so if the cluster is simply for the purposes of running this demo, it is recommended to only use a single manager node).

```console
$ docker stack deploy -c stack.yml --resolve-image never kubecon-demo
```

Once running, hit `http://localhost:9090` or `http://swarm-ip:9090` in your web browser and you should see a list of all worker nodes and their associated architecture, with a button to toggle their involvement in the network generation.
