package hoard

type MigrationOptions struct {
	Force bool
}

type MirrorRepoOptions struct {
	MigrationOptions

	OwnerName      string
	RepositoryName string
}

type MirrorProfileOptions struct {
	MigrationOptions

	OwnerName string
}
