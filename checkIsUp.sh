#!/bin/sh

#Check if server is up

if [[ -z "${SCISCRAPER_CONDUCTOR_HOST}" ]]; then
  echo ERRO: SCISCRAPER_CONDUCTOR_HOST não informado
  exit 1
fi

if [[ -z "${SCISCRAPER_CONDUCTOR_PORT}" ]]; then
  echo ERRO: SCISCRAPER_CONDUCTOR__PORT não informado
  exit 1
fi

while [ "$(nc -z $SCISCRAPER_CONDUCTOR_HOST $SCISCRAPER_CONDUCTOR_PORT </dev/null; echo $?)" !=  "0" ];
do sleep 5;
echo "Waiting until CONDUCTOR SERVER is UP and RESPONDING";
done

sleep 5;

./scicrop-scraper run
