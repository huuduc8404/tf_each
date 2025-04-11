resource "aws_instance" "web1" {
  ami           = "ami-123456"
  instance_type = "t2.micro"
  tags = {
    Name = "WebServer1"
  }
}

resource "aws_instance" "web2" {
  ami           = "ami-789012"
  instance_type = "t2.small"
  tags = {
    Name = "WebServer2"
  }
}
