build:
	@go build -o ./bin/forres ./cmd/api

run: build
	@./bin/forres

test:
	@go test ./cmd/test


# User
# 	- Sign in 
# 	- Sign up
# 	=> User Tables

# Customer
# 	=> Voucher
# 	=> Course
# 	=> Menu Course
# 	=> Customer 
# 	=> Booking

# Staff
# 	=> Booking