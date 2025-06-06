package setting

import (
	"context"
	"github.com/NYCU-SDC/summer/pkg/database"
	"github.com/NYCU-SDC/summer/pkg/log"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Service struct {
	logger *zap.Logger
	tracer trace.Tracer
	query  *Queries
}

func NewService(logger *zap.Logger, db DBTX) *Service {
	return &Service{
		logger: logger,
		tracer: otel.Tracer("setting/service"),
		query:  New(db),
	}
}

func (s *Service) GetSettingByUserID(ctx context.Context, userID uuid.UUID) (Setting, error) {
	traceCtx, span := s.tracer.Start(ctx, "GetSettingByUserID")
	defer span.End()
	logger := logutil.WithContext(traceCtx, s.logger)

	setting, err := s.query.GetSetting(ctx, userID)
	if err != nil {
		err = databaseutil.WrapDBErrorWithKeyValue(err, "settings", "id", userID.String(), logger, "get setting by user id")
		span.RecordError(err)
		return Setting{}, err
	}

	return setting, nil
}

func (s *Service) FindOrCreateSetting(ctx context.Context, userID uuid.UUID, username pgtype.Text) (Setting, error) {
	traceCtx, span := s.tracer.Start(ctx, "UpdateSetting")
	defer span.End()
	logger := logutil.WithContext(traceCtx, s.logger)

	exist, err := s.query.SettingExists(ctx, userID)
	if err != nil {
		err = databaseutil.WrapDBErrorWithKeyValue(err, "settings", "id", userID.String(), logger, "check setting exists")
		span.RecordError(err)
		return Setting{}, err
	}

	var setting Setting
	if !exist {
		setting, err = s.query.CreateSetting(ctx, CreateSettingParams{UserID: userID, Username: username})
		if err != nil {
			err = databaseutil.WrapDBError(err, logger, "create setting")
			span.RecordError(err)
			return Setting{}, err
		}
	} else {
		setting, err = s.query.GetSetting(ctx, userID)
		if err != nil {
			err = databaseutil.WrapDBErrorWithKeyValue(err, "settings", "id", userID.String(), logger, "get setting by user id")
			span.RecordError(err)
			return Setting{}, err
		}
	}
	
	return setting, nil
}

func (s *Service) UpdateSetting(ctx context.Context, userID uuid.UUID, setting Setting) (Setting, error) {
	traceCtx, span := s.tracer.Start(ctx, "UpdateSetting")
	defer span.End()
	logger := logutil.WithContext(traceCtx, s.logger)

	updatedSetting, err := s.query.UpdateSetting(ctx, UpdateSettingParams(setting))
	if err != nil {
		err = databaseutil.WrapDBErrorWithKeyValue(err, "settings", "id", userID.String(), logger, "update setting")
		span.RecordError(err)
		return Setting{}, err
	}

	return updatedSetting, nil
}

func (s *Service) GetPublicKeysByUserID(ctx context.Context, userID uuid.UUID) ([]PublicKey, error) {
	traceCtx, span := s.tracer.Start(ctx, "GetPublicKeysByUserID")
	defer span.End()
	logger := logutil.WithContext(traceCtx, s.logger)

	publicKeys, err := s.query.GetPublicKeys(ctx, userID)
	if err != nil {
		err = databaseutil.WrapDBErrorWithKeyValue(err, "public_keys", "id", userID.String(), logger, "get public keys by user id")
		span.RecordError(err)
		return nil, err
	}

	return publicKeys, err
}

func (s *Service) GetPublicKeyByID(ctx context.Context, id uuid.UUID) (PublicKey, error) {
	traceCtx, span := s.tracer.Start(ctx, "GetPublicKeyByID")
	defer span.End()
	logger := logutil.WithContext(traceCtx, s.logger)

	publicKey, err := s.query.GetPublicKey(ctx, id)
	if err != nil {
		err = databaseutil.WrapDBErrorWithKeyValue(err, "public_keys", "id", id.String(), logger, "get public key by id")
		span.RecordError(err)
		return PublicKey{}, err
	}

	return publicKey, nil
}

func (s *Service) AddPublicKey(ctx context.Context, publicKey AddPublicKeyParams) (PublicKey, error) {
	traceCtx, span := s.tracer.Start(ctx, "AddPublicKey")
	defer span.End()
	logger := logutil.WithContext(traceCtx, s.logger)

	addedPublicKey, err := s.query.AddPublicKey(ctx, publicKey)
	if err != nil {
		err = databaseutil.WrapDBErrorWithKeyValue(err, "public_keys", "id", publicKey.UserID.String(), logger, "add public key")
		span.RecordError(err)
		return PublicKey{}, err
	}

	return addedPublicKey, nil
}

func (s *Service) DeletePublicKey(ctx context.Context, id uuid.UUID) error {
	traceCtx, span := s.tracer.Start(ctx, "DeletePublicKey")
	defer span.End()
	logger := logutil.WithContext(traceCtx, s.logger)

	err := s.query.DeletePublicKey(ctx, id)
	if err != nil {
		err = databaseutil.WrapDBErrorWithKeyValue(err, "settings", "id", id.String(), logger, "delete public key")
		span.RecordError(err)
		return err
	}

	return nil
}
