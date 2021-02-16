package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/llamadeus/iot-logger/features/channels"
	"github.com/llamadeus/iot-logger/graph/generated"
	"github.com/llamadeus/iot-logger/graph/types"
)

func (r *queryResolver) Ping(ctx context.Context) (string, error) {
	return "pong", nil
}

func (r *queryResolver) History(ctx context.Context, channel string) ([]*types.Message, error) {
	return channels.History(ctx, channel)
}

func (r *mutationResolver) AddMessage(ctx context.Context, channel string, message string) (bool, error) {
	return channels.AddMessage(ctx, channel, message)
}

func (r *subscriptionResolver) MessageAdded(ctx context.Context, channel string) (<-chan *types.Message, error) {
	return channels.MessageAdded(ctx, channel)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
