#!/bin/sh

# Generate mocks for repository interfaces
mockery --name=CommonRepository --dir=internal/repository/common --output=internal/test/mockrepository --outpkg=mockrepository
mockery --name=UserRepository --dir=internal/repository/user --output=internal/test/mockrepository --outpkg=mockrepository

# Generate mocks for service interfaces
mockery --name=AuthService --dir=internal/service/auth --output=internal/test/mockservice --outpkg=mockservice