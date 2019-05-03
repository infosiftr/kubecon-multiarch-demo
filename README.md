# Running the Demo

Set up a cluster of `amd64`, `arm32v7`, and/or `arm64v8` machines running kubernetes.  The application will deploy as a daemon set, based on the label defined in `deploy.yml`.  The application will run on any node labeled with `app=blinky-node` and will cease running when relabeled with `app=blinky-nope`

First, deploy the web server locally

```console
$ ./run-server.sh
```

Then deploy the application across the kubernetes cluster

```console
$ kubectl apply -f deploy.yml
```

Once running, hit `http://localhost:8080` in your web browser and you should see a list of all worker nodes and their associated architecture, with a button to toggle their involvement in the network generation.
