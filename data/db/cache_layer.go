package db

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/joshuaschlichting/gocms/data/cache"
)

type DBCache struct {
	queries *Queries
	cache   *cache.Cache
}

func NewDBCache(q *Queries, c *cache.Cache) *DBCache {
	return &DBCache{queries: q, cache: c}
}

// Wrapper for GetBlogPost with caching
func (c *DBCache) GetBlogPost(ctx context.Context, id uuid.UUID) (GetBlogPostRow, error) {
	cacheKey := "BlogPost:" + id.String()
	log.Println("Getting blog post from cache")
	if val, err := c.cache.Get(cacheKey); err == nil {
		if bp, ok := val.(GetBlogPostRow); ok {
			return bp, nil
		}
	}

	bp, err := c.queries.GetBlogPost(ctx, id)
	if err == nil {
		c.cache.Set(cacheKey, bp, time.Minute*5)
	}
	return bp, err
}

// Wrapper for GetUser with caching
func (c *DBCache) GetUser(ctx context.Context, id uuid.UUID) (User, error) {
	cacheKey := "User:" + id.String()

	if val, err := c.cache.Get(cacheKey); err == nil {
		if user, ok := val.(User); ok {
			log.Println("Got user from cache")
			return user, nil
		}
	}

	user, err := c.queries.GetUser(ctx, id)
	if err == nil {
		c.cache.Set(cacheKey, user, time.Minute*5)
	}
	return user, err
}

// Wrapper for CreateBlogPost with caching
func (c *DBCache) CreateBlogPost(ctx context.Context, arg CreateBlogPostParams) (BlogPost, error) {
	bp, err := c.queries.CreateBlogPost(ctx, arg)
	if err == nil {
		cacheKey := "BlogPost:" + bp.ID.String()
		c.cache.Set(cacheKey, bp, time.Minute*5) // Cache for 5 minutes
	}
	return bp, err
}

// Wrapper for CreateMessage (though for creation, it's mostly just passing through)
func (c *DBCache) CreateMessage(ctx context.Context, arg CreateMessageParams) (Message, error) {
	// For creations, typically we just call the DB without caching the result.
	// But you might want to cache specific queries if needed.
	return c.queries.CreateMessage(ctx, arg)
}

// Wrapper for CreateOrganization with caching
func (c *DBCache) CreateOrganization(ctx context.Context, arg CreateOrganizationParams) (Organization, error) {
	org, err := c.queries.CreateOrganization(ctx, arg)
	if err == nil {
		cacheKey := "Organization:" + org.ID.String()
		c.cache.Set(cacheKey, org, time.Minute*5)
	}
	return org, err
}

// Wrapper for GetUserByName with caching
func (c *DBCache) GetUserByName(ctx context.Context, name string) (User, error) {
	cacheKey := "UserByName:" + name

	if val, err := c.cache.Get(cacheKey); err == nil {
		if user, ok := val.(User); ok {
			return user, nil
		}
	}

	user, err := c.queries.GetUserByName(ctx, name)
	if err == nil {
		c.cache.Set(cacheKey, user, time.Minute*5)
	}
	return user, err
}

// Wrapper for GetUserIsInGroup with caching
func (c *DBCache) GetUserIsInGroup(ctx context.Context, arg GetUserIsInGroupParams) (bool, error) {
	cacheKey := "UserIsInGroup:" + arg.UserID.String() + ":" + arg.UsergroupName

	if val, err := c.cache.Get(cacheKey); err == nil {
		if isInGroup, ok := val.(bool); ok {
			return isInGroup, nil
		}
	}

	isInGroup, err := c.queries.GetUserIsInGroup(ctx, arg)
	if err == nil {
		c.cache.Set(cacheKey, isInGroup, time.Minute*5)
	}
	return isInGroup, err
}

// Wrapper for ListBlogPosts with caching
func (c *DBCache) ListBlogPosts(ctx context.Context) ([]ListBlogPostsRow, error) {
	cacheKey := "ListBlogPosts"

	if val, err := c.cache.Get(cacheKey); err == nil {
		if bpList, ok := val.([]ListBlogPostsRow); ok {
			return bpList, nil
		}
	}

	bpList, err := c.queries.ListBlogPosts(ctx)
	if err == nil {
		c.cache.Set(cacheKey, bpList, time.Minute*5)
	}
	return bpList, err
}

// Wrapper for ListBlogPostsByUser with caching
func (c *DBCache) ListBlogPostsByUser(ctx context.Context, authorID uuid.UUID) ([]ListBlogPostsByUserRow, error) {
	cacheKey := "ListBlogPostsByUser:" + authorID.String()

	if val, err := c.cache.Get(cacheKey); err == nil {
		if bpList, ok := val.([]ListBlogPostsByUserRow); ok {
			return bpList, nil
		}
	}

	bpList, err := c.queries.ListBlogPostsByUser(ctx, authorID)
	if err == nil {
		c.cache.Set(cacheKey, bpList, time.Minute*5)
	}
	return bpList, err
}

// Wrapper for ListFiles with caching
func (c *DBCache) ListFiles(ctx context.Context) ([]File, error) {
	cacheKey := "ListFiles"

	if val, err := c.cache.Get(cacheKey); err == nil {
		if fileList, ok := val.([]File); ok {
			return fileList, nil
		}
	}

	fileList, err := c.queries.ListFiles(ctx)
	if err == nil {
		c.cache.Set(cacheKey, fileList, time.Minute*5)
	}
	return fileList, err
}

// Wrapper for ListMessagesFrom with caching
func (c *DBCache) ListMessagesFrom(ctx context.Context, fromID uuid.UUID) ([]ListMessagesFromRow, error) {
	cacheKey := "ListMessagesFrom:" + fromID.String()

	if val, err := c.cache.Get(cacheKey); err == nil {
		if msgList, ok := val.([]ListMessagesFromRow); ok {
			return msgList, nil
		}
	}

	msgList, err := c.queries.ListMessagesFrom(ctx, fromID)
	if err == nil {
		c.cache.Set(cacheKey, msgList, time.Minute*5)
	}
	return msgList, err
}

// Wrapper for ListMessagesTo with caching
func (c *DBCache) ListMessagesTo(ctx context.Context, toUsername string) ([]ListMessagesToRow, error) {
	cacheKey := "ListMessagesTo:" + toUsername

	if val, err := c.cache.Get(cacheKey); err == nil {
		if msgList, ok := val.([]ListMessagesToRow); ok {
			return msgList, nil
		}
	}

	msgList, err := c.queries.ListMessagesTo(ctx, toUsername)
	if err == nil {
		c.cache.Set(cacheKey, msgList, time.Minute*5)
	}
	return msgList, err
}

// Wrapper for CreateUser without caching but invalidates relevant cache
func (c *DBCache) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	user, err := c.queries.CreateUser(ctx, arg)
	if err == nil {
		cacheKey := "User:" + user.ID.String()
		c.cache.Delete(cacheKey)
	}
	return user, err
}

// Wrapper for UpdateUser without caching but invalidates relevant cache
func (c *DBCache) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	user, err := c.queries.UpdateUser(ctx, arg)
	if err == nil {
		cacheKey := "User:" + user.ID.String()
		c.cache.Delete(cacheKey)
	}
	return user, err
}

// Wrapper for DeleteUser without caching but invalidates relevant cache
func (c *DBCache) DeleteUser(ctx context.Context, id uuid.UUID) error {
	err := c.queries.DeleteUser(ctx, id)
	if err == nil {
		cacheKey := "User:" + id.String()
		c.cache.Delete(cacheKey)
	}
	return err
}

// Wrapper for UpdateOrganization without caching but invalidates relevant cache
func (c *DBCache) UpdateOrganization(ctx context.Context, arg UpdateOrganizationParams) (Organization, error) {
	org, err := c.queries.UpdateOrganization(ctx, arg)
	if err == nil {
		cacheKey := "Organization:" + org.ID.String()
		c.cache.Delete(cacheKey)
	}
	return org, err
}

// Wrapper for DeleteOrganization without caching but invalidates relevant cache
func (c *DBCache) DeleteOrganization(ctx context.Context, id uuid.UUID) error {
	err := c.queries.DeleteOrganization(ctx, id)
	if err == nil {
		cacheKey := "Organization:" + id.String()
		c.cache.Delete(cacheKey)
	}
	return err
}

// Wrapper for CreateUserGroup without caching but invalidates relevant cache
func (c *DBCache) CreateUserGroup(ctx context.Context, arg CreateUserGroupParams) (Usergroup, error) {
	ug, err := c.queries.CreateUserGroup(ctx, arg)
	if err == nil {
		cacheKey := "UserGroup:" + ug.ID.String()
		c.cache.Delete(cacheKey)
	}
	return ug, err
}

// Wrapper for UpdateUserGroup without caching but invalidates relevant cache
func (c *DBCache) UpdateUserGroup(ctx context.Context, arg UpdateUserGroupParams) (Usergroup, error) {
	ug, err := c.queries.UpdateUserGroup(ctx, arg)
	if err == nil {
		cacheKey := "UserGroup:" + ug.ID.String()
		c.cache.Delete(cacheKey)
	}
	return ug, err
}

// Wrapper for DeleteUserGroup without caching but invalidates relevant cache
func (c *DBCache) DeleteUserGroup(ctx context.Context, id uuid.UUID) error {
	err := c.queries.DeleteUserGroup(ctx, id)
	if err == nil {
		cacheKey := "UserGroup:" + id.String()
		c.cache.Delete(cacheKey)
	}
	return err
}

func (c *DBCache) ListOrganizations(ctx context.Context) ([]Organization, error) {
	cacheKey := "OrganizationsList"

	if val, err := c.cache.Get(cacheKey); err == nil {
		if orgs, ok := val.([]Organization); ok {
			return orgs, nil
		}
	}

	orgs, err := c.queries.ListOrganizations(ctx)
	if err == nil {
		c.cache.Set(cacheKey, orgs, time.Minute*5)
	}
	return orgs, err
}

func (c *DBCache) ListUsers(ctx context.Context) ([]User, error) {
	cacheKey := "UsersList"

	if val, err := c.cache.Get(cacheKey); err == nil {
		if users, ok := val.([]User); ok {
			return users, nil
		}
	}

	users, err := c.queries.ListUsers(ctx)
	if err == nil {
		c.cache.Set(cacheKey, users, time.Minute*5)
	}
	return users, err
}

func (c *DBCache) ListUserGroups(ctx context.Context) ([]Usergroup, error) {
	cacheKey := "UserGroupsList"

	if val, err := c.cache.Get(cacheKey); err == nil {
		if userGroups, ok := val.([]Usergroup); ok {
			return userGroups, nil
		}
	}

	userGroups, err := c.queries.ListUserGroups(ctx)
	if err == nil {
		c.cache.Set(cacheKey, userGroups, time.Minute*5)
	}
	return userGroups, err
}
