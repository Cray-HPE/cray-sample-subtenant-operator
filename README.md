# cray-sample-subtenant-operator (Sample implementation of sub-tenant operator)

## Overview

This repo is intended simply to provide an example of a kubernetes operator that watches the [cray-tapms-operator](https://github.com/Cray-HPE/cray-tapms-operator) CRD and can respond to changes accordingly.  This controller follows the Operator Pattern (https://kubernetes.io/docs/concepts/extend-kubernetes/operator/).

## Notes

See [subtenant_controller.go](https://github.com/Cray-HPE/cray-sample-subtenant-operator/blob/main/controllers/subtenant_controller.go) for the code that watches and handles the tapms CRD.

