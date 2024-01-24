FROM fedora:39

# Install go
RUN dnf update -y
RUN dnf install go -y

# Copy and Build
WORKDIR /app
COPY src/zl16 .
RUN go build -o zl16
RUN rm -rf banner.txt go.mod go.sum main.go names.txt http/custom-contacts.csv http/contacts.csv

# Expost port
EXPOSE 1999

# Run app
CMD ./zl16

