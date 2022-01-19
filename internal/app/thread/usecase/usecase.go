package thread_usecase

import (
	"strconv"

	"github.com/perlinleo/technopark-mail.ru-forum-database/internal/app/thread"
	"github.com/perlinleo/technopark-mail.ru-forum-database/internal/app/user"
	"github.com/perlinleo/technopark-mail.ru-forum-database/internal/model"
	"github.com/pkg/errors"
)

type ThreadUsecase struct {
	threadRep thread.Repository
	userRep   user.Repository
}

func (t ThreadUsecase) Vote(threadSlugOrId string, vote *model.Vote) (*model.Thread, error) {
	threadObj, err := t.FindByIdOrSlug(threadSlugOrId)
	
	if err != nil {
		return nil, err
	}

	threadObj, err = t.threadRep.Vote(threadObj, vote)
	if err != nil {
		return nil, err
	}

	return threadObj, nil
}

func (t ThreadUsecase) GetThreadPosts(threadSlugOrId string, limit string, desc bool, since string, sort string) ([]model.Post, error) {
	threadObj, err := t.FindByIdOrSlug(threadSlugOrId) 
	if err != nil {
		return nil, err
	}

	if limit == "" {
		limit = "100"
	}
	descValue :=  ""
	if desc==true {
		descValue = "desc"
	}

	// mb cringe 
	
	if sort=="" {
		sort = "flat"
	}
	

	posts, err := t.threadRep.GetThreadPosts(threadObj, limit, descValue, since, sort)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (t ThreadUsecase) UpdateThread(threadSlugOrId string, threadUpdate *model.ThreadUpdate) (*model.Thread, error) {
	id, _ := strconv.Atoi(threadSlugOrId)

	threadObj, err := t.threadRep.FindByIdOrSlug(id, threadSlugOrId)
	if err != nil {
		return nil, err
	}

	if threadUpdate.Title=="" && threadUpdate.Message=="" {
		return threadObj, nil
	}

	if threadUpdate.Title=="" {
		threadUpdate.Title = threadObj.Title
	}

	if threadUpdate.Message=="" {
		threadUpdate.Message = threadObj.Message
	}

	threadObj, err = t.threadRep.UpdateThread(id, threadSlugOrId, threadUpdate)
	if err != nil {
		return nil, err
	}

	return threadObj, nil
}

func (t ThreadUsecase) FindByIdOrSlug(threadSlugOrId string) (*model.Thread, error) {
	id, _ := strconv.Atoi(threadSlugOrId)

	threadObj, err := t.threadRep.FindByIdOrSlug(id, threadSlugOrId)
	if threadObj == nil || err != nil {
		return nil, err
	}

	return threadObj, nil
}

func (t ThreadUsecase) CreatePosts(threadSlugOrId string, posts []*model.Post) ([]*model.Post, int, error) {
	threadObj, err := t.FindByIdOrSlug(threadSlugOrId)
	if threadObj == nil || err != nil {
		return nil, 404, err
	}

	posts, err = t.threadRep.CreatePosts(threadObj, posts)
	if err != nil {
		if err.Error() == "404" {
			return nil, 404, err
		}

		return nil, 409, err
	}

	return posts, 201, nil
}

func NewThreadUsecase(t thread.Repository, u user.Repository) thread.Usecase {
	return &ThreadUsecase{
		threadRep: t,
		userRep:   u,
	}
}

func makeTree(posts []model.Post) []model.Post {
	tree := make([]model.Post, 0)
	var parent *model.Post

	for _, p := range posts {
		if len(p.Path) == 1 {
			tree = append(tree, p)
			parent = &p
		} else if len(p.Path) > 1 {
			if p.Parent == parent.ID {
				//parent.Childs = append(parent.Childs, p)
				tree = append(tree, p)
				p.ParentPointer = parent
				parent = &p
			} else {
				for p.Parent != parent.ID {
					parent = parent.ParentPointer
				}
				//parent.Childs = append(parent.Childs, p)
				tree = append(tree, p)
				p.ParentPointer = parent
				parent = &p
			}
		}
	}

	return tree
}


func (t ThreadUsecase) FindPostId(id string, includeUser bool,includeForum bool, includeThread bool) (*model.PostFull, error) {
	
	

	postObj, err := t.threadRep.FindPostId(id, includeUser, includeForum, includeThread)

	if err != nil {
		return nil, errors.Wrap(err, "postRep.FindById()")
	}

	return postObj, nil
}


func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}

	return false
}


func (t ThreadUsecase) UpdatePost(id string, message string) (*model.Post, error) {
	postFullObj, err := t.threadRep.FindPostId(id, false, false, false)
	if err != nil {
		return nil, errors.Wrap(err, "postRep.FindById()")
	}

	if message=="" || postFullObj.Post.Message == message {
		return postFullObj.Post, nil
	}

	postObj, err := t.threadRep.UpdatePost(id, message)

	if err != nil {
		return nil, errors.Wrap(err, "postRep.Update()")
	}

	return postObj, nil
}


func (t ThreadUsecase) ClearAll() error {
	err := t.threadRep.ClearAll()

	if err != nil {
		return errors.Wrap(err, "ClearAll()")
	}

	return nil
}

func (t ThreadUsecase) GetStatus() (*model.Status, error) {
	status, err := t.threadRep.GetStatus()

	if err != nil {
		return nil, errors.Wrap(err, "GetStatus()")
	}

	return status, nil
}