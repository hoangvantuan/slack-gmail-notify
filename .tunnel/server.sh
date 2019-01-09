#!/bin/bash

tunneld -tlsCrt /root/.tunneld/fullchain.pem -tlsKey /root/.tunneld/privkey.pem > /dev/null 2>&1 &