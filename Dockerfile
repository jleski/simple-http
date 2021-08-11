FROM scratch
# Copy our static executable.
COPY ./main /main
# Run the binary.
ENTRYPOINT ["/main"]