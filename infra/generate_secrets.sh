ENV=$1

case $ENV in
    "dev")
        echo $ENV
        cat ./secrets/$ENV/backend.yaml | kubeseal \
            --controller-namespace kube-system \
            --controller-name sealed-secrets-controller \
            --format json \
            > ./helm/templates/sealed-secret.json
        ;;
    *)
        echo "invalid env (got '$ENV', expected 'prod'/'dev')"
        exit 1
        ;;
esac
