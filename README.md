# Varnish Stats to Google Monitoring v3

Supports collecting the following Varnish counters and pushing to Google Cloud Monitoring v3 every minute.

Uses

```
"MAIN.sess_conn", "MAIN.sess_drop", "MAIN.sess_fail", "MAIN.cache_hit",
"MAIN.cache_hitpass", "MAIN.cache_miss", "MAIN.client_req", "MAIN.backend_conn",
"MAIN.backend_busy", "MAIN.backend_reuse", "MAIN.threads", "MAIN.n_object",
"MAIN.n_lru_nuked", "MAIN.bans_obj", "MAIN.bans_req"
```
## Running
Collection loop is timed to every minute.  Metric value pushed is a delta calculation from prior minute and current minute.  Google Metrics are created as gauge's.

Uses computeMetadata service to pull projectID, instance name, and the Custom Metadata "artifact-id" to use as application label on metric.

Uses a blocking ticker loop to "daemonize", but needs an init script.  Make sure the user running gogo-varnishstat has access to the VSM.  There's no config / switches needed.

Only supports the default VSM name (no `-n someapp` support).  Basically, if the user can `varnishstat -1`, gogo-varnishstat should work.

Requires Monitoring V3 Write scope.  

## Building

Uses github.com/phenomenes/vago for varnish VSM access.  Vago has a few external build reqs.  It requires `pkg-config` and `libvarnishapi-dev` (or `varnish-plus-dev` for varnish plus).  Additionally an env var should be present for the `go build` which points to the location of the `varnishapi.pc`.  The following works for me on Ubuntu 14/16, but YMMV.  
```
PKG_CONFIG_PATH=/usr/lib/pkgconfig
```
See https://github.com/phenomenes/vago for more build details.  The rest of the code is standard Go.
