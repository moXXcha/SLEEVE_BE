package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		// 内部ID（auto increment、外部には公開しない）
		// entのデフォルトIDはintのauto incrementなので、明示的な定義は不要
		// 公開ID（UUID、外部に公開する）
		field.UUID("public_id", uuid.UUID{}).
			Default(uuid.New).
			Unique().
			Immutable().
			Comment("公開用ユーザーID（UUID）"),
		field.String("firebase_uid").
			NotEmpty().
			Unique().
			Comment("Firebase Authentication UID"),
		field.String("email").
			NotEmpty().
			Unique().
			Comment("メールアドレス"),
		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Comment("作成日時"),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("更新日時"),
		field.Time("deleted_at").
			Optional().
			Nillable().
			Comment("削除日時（論理削除）"),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}

// Indexes of the User.
func (User) Indexes() []ent.Index {
	return []ent.Index{
		// public_idにインデックス（外部からの検索用）
		index.Fields("public_id").
			Unique(),
		// firebase_uidにインデックス
		index.Fields("firebase_uid").
			Unique(),
		// emailにインデックス
		index.Fields("email").
			Unique(),
		// deleted_atにインデックス（論理削除の検索用）
		index.Fields("deleted_at"),
	}
}
