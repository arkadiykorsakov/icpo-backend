package dto

import tokenDto "icpo-backend/internal/domain/token/dto"

func ToTokenPairResponse(tokenPair tokenDto.TokenPair) TokenPairResponse {
	return TokenPairResponse{
		AccessToken: tokenPair.AccessToken,
		TokenType:   tokenPair.TokenType,
		ExpiresIn:   tokenPair.ExpiresIn,
	}
}
