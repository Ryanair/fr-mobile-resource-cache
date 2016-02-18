# Recource caching using Sync Gateway

![Build status](https://travis-ci.org/Ryanair/fr-mobile-resource-cache.svg?branch=master)

The service allows you to maintain a cache of json and binary documents which can be sync-ed to remote clients, using [**Couchbase Sync gateway**](https://github.com/couchbase/sync_gateway) and [**Couchbase lite**](http://developer.couchbase.com/mobile/). This technique is useful for distributing semi-static content and images. The client reads the data locally and the couchbase lite client and sync gateway hide the syncronisation complexity from the developer. 

The service listens for changes in a directory and generates new revisions, which triggers sync to the remote clients. This way there is no wasted bandwith between the remote client and our server and applications don't have to wait to serve content to the user.

You can find a demo project in this [**repository**](https://github.com/Ryanair/resource-sync-example/tree/master).

![](http://i284.photobucket.com/albums/ll17/Vlado_Atanasov/go_resource_update_zps61xhoepx.png)

![](http://i284.photobucket.com/albums/ll17/Vlado_Atanasov/animation_zpssxqbookb.gif)
