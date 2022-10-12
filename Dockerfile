FROM gcr.io/distroless/static:nonroot
COPY kube-webhook-certgen /kube-webhook-certgen
USER 65532:65532
ENTRYPOINT ["/kube-webhook-certgen"]
