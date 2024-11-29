#!/bin/sh

# Generate mocks for repository interfaces
mockery --name=CommonRepository --dir=internal/repository/common --output=internal/test/mockrepository --outpkg=mockrepository
mockery --name=UserRepository --dir=internal/repository/user --output=internal/test/mockrepository --outpkg=mockrepository
mockery --name=UserMatchRepository --dir=internal/repository/user_match --output=internal/test/mockrepository --outpkg=mockrepository
mockery --name=UserPremiumRepository --dir=internal/repository/user_premium --output=internal/test/mockrepository --outpkg=mockrepository
mockery --name=PremiumConfigRepository --dir=internal/repository/premium_config --output=internal/test/mockrepository --outpkg=mockrepository

# Generate mocks for service interfaces
mockery --name=AuthService --dir=internal/service/auth --output=internal/test/mockservice --outpkg=mockservice

# Generate mocks for usecase interfaces
mockery --name=UserUsecase --dir=internal/usecase/user --output=internal/test/mockusecase --outpkg=mockusecase
mockery --name=PremiumConfigUsecase --dir=internal/usecase/premium_config --output=internal/test/mockusecase --outpkg=mockusecase
mockery --name=UserMatchUsecase --dir=internal/usecase/user_match --output=internal/test/mockusecase --outpkg=mockusecase