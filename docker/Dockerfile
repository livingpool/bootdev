# This is a comment

# Use a lightweight debian os
# as the base image
FROM debian:stable-slim

# COPY source destination
# It's also worth mentioning that the ADD command would also work here, 
# but COPY is the better choice for us because we don't need the extra functionality that ADD offers.
COPY docker /bin/docker

ENV PORT=8080

# command when the container runs
CMD ["/bin/docker"]
