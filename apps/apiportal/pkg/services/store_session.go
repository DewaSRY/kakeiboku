package services

import (
	"context"
	"database/sql"
)


func (s *SQLStore) SetSession(ctx context.Context, arg CreateSessionParams) (Session, error) {

	current_session, err := s.GetSessionsByEmailAndIp(ctx, GetSessionsByEmailAndIpParams{
		Email:    arg.Email,
		ClientIp: arg.ClientIp,
		UserAgent: arg.UserAgent,
	})

	if err != nil {
		if err == sql.ErrNoRows {
			session, err := s.CreateSession(ctx, arg)
			if err != nil {
				return Session{}, err
			}
			return session, nil
		}
	}

	update_session, err := s.RefreshSession(ctx, RefreshSessionParams{
		ID: current_session.ID,
		RefreshToken: arg.RefreshToken,
		ExpiresAt: arg.ExpiresAt,
	})

	if err != nil {
		return Session{}, err
	}
	
	return update_session, nil
}	