# prometheus2tower

Prometheus2tower is a service that integrates [prometheus-alertmanager](https://prometheus.io/docs/alerting/latest/alertmanager/) and [Ansible Tower](https://docs.ansible.com/ansible-tower/).

Its primary use is to trigger playbook runs (e.g. remediation playbooks) in Tower on specific alerts being fired by prometheus.

It listens on [configured](cmd/prom2tower/conf.yaml.example) endpoints for incoming alerts. Upon receipt, unmarshalls received JSON and can use this data to fill out the template of the configured tower request body.
The request is then sent to the configured tower endpoint to trigger a template run, or, it might be configured to do other [Tower API](https://docs.ansible.com/ansible-tower/latest/html/towerapi/index.html) calls.

Example endpoint configuration:

```
  - name: alertmanager-tower-glue
    ingress:
      path: "/webhook-url-configured-as-alertmanager-receiver"
    egress:
      towerHost: "https://tower.local.lan"
      towerToken: "PasteLocalTokenHere"
      path: "/api/v2/job_templates/364/launch/"
      method: "POST"
      body: >
        { 
          "limit": "{{ range .Alerts }}{{ .Labels.instance }},{{ end }}",
          "verbosity": 3,
          "extra_vars": 
            { 
              "input": "Firing for: {{ range .Alerts }}{{ .Labels.instance }} {{ end }}"
            } 
        }
```

## notes:

* tower: json 
    * e.g. https://tower/api/v2/job_templates/364/launch/
    * NOTE:
    * Every field sent in json body MUST have "Prompt on launch" checkbox enabled.
    * For example, to enable sending "limit", template must have: `"ask_limit_on_launch": true`


---

**Disclaimer Warning**: this is under development and there are no guarantees that it works correctly... or at all.

