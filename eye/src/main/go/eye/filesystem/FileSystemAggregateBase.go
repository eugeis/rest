package filesystem

import (
    "errors"
    "fmt"
    "github.com/eugeis/gee/eh"
    "github.com/looplab/eventhorizon"
    "github.com/looplab/eventhorizon/commandhandler/bus"
    "time"
)
type CommandHandler struct {
    CreateHandler func (*Create, *FileSystemService, eh.AggregateStoreEvent) (err error)  `json:"createHandler" eh:"optional"`
    DeleteHandler func (*Delete, *FileSystemService, eh.AggregateStoreEvent) (err error)  `json:"deleteHandler" eh:"optional"`
    UpdateHandler func (*Update, *FileSystemService, eh.AggregateStoreEvent) (err error)  `json:"updateHandler" eh:"optional"`
}

func (o *CommandHandler) AddCreatePreparer(preparer func (*Create, *FileSystemService) (err error) ) {
    prevHandler := o.CreateHandler
	o.CreateHandler = func(command *Create, entity *FileSystemService, store eh.AggregateStoreEvent) (err error) {
		if err = preparer(command, entity); err == nil {
			err = prevHandler(command, entity, store)
		}
		return
	}
}

func (o *CommandHandler) AddDeletePreparer(preparer func (*Delete, *FileSystemService) (err error) ) {
    prevHandler := o.DeleteHandler
	o.DeleteHandler = func(command *Delete, entity *FileSystemService, store eh.AggregateStoreEvent) (err error) {
		if err = preparer(command, entity); err == nil {
			err = prevHandler(command, entity, store)
		}
		return
	}
}

func (o *CommandHandler) AddUpdatePreparer(preparer func (*Update, *FileSystemService) (err error) ) {
    prevHandler := o.UpdateHandler
	o.UpdateHandler = func(command *Update, entity *FileSystemService, store eh.AggregateStoreEvent) (err error) {
		if err = preparer(command, entity); err == nil {
			err = prevHandler(command, entity, store)
		}
		return
	}
}

func (o *CommandHandler) Execute(cmd eventhorizon.Command, entity eventhorizon.Entity, store eh.AggregateStoreEvent) (err error) {
    switch cmd.CommandType() {
    case CreateCommand:
        err = o.CreateHandler(cmd.(*Create), entity.(*FileSystemService), store)
    case DeleteCommand:
        err = o.DeleteHandler(cmd.(*Delete), entity.(*FileSystemService), store)
    case UpdateCommand:
        err = o.UpdateHandler(cmd.(*Update), entity.(*FileSystemService), store)
    default:
		err = errors.New(fmt.Sprintf("Not supported command type '%v' for entity '%v", cmd.CommandType(), entity))
	}
    return
}

func (o *CommandHandler) SetupCommandHandler() (err error) {
    o.CreateHandler = func(command *Create, entity *FileSystemService, store eh.AggregateStoreEvent) (err error) {
        if err = eh.ValidateNewId(entity.Id, command.Id, FileSystemServiceAggregateType); err == nil {
            store.StoreEvent(createdEvent, &Created{
                Name: command.Name,
                Id: command.Id,}, time.Now())
        }
        return
    }
    o.DeleteHandler = func(command *Delete, entity *FileSystemService, store eh.AggregateStoreEvent) (err error) {
        if err = eh.ValidateIdsMatch(entity.Id, command.Id, FileSystemServiceAggregateType); err == nil {
            store.StoreEvent(deletedEvent, &Deleted{
                Id: command.Id,}, time.Now())
        }
        return
    }
    o.UpdateHandler = func(command *Update, entity *FileSystemService, store eh.AggregateStoreEvent) (err error) {
        if err = eh.ValidateIdsMatch(entity.Id, command.Id, FileSystemServiceAggregateType); err == nil {
            store.StoreEvent(updatedEvent, &Updated{
                Name: command.Name,
                Id: command.Id,}, time.Now())
        }
        return
    }
    return
}


type EventHandler struct {
    CreateHandler func (*Create, *FileSystemService) (err error)  `json:"createHandler" eh:"optional"`
    CreatedHandler func (*Created, *FileSystemService) (err error)  `json:"createdHandler" eh:"optional"`
    DeleteHandler func (*Delete, *FileSystemService) (err error)  `json:"deleteHandler" eh:"optional"`
    DeletedHandler func (*Deleted, *FileSystemService) (err error)  `json:"deletedHandler" eh:"optional"`
    UpdateHandler func (*Update, *FileSystemService) (err error)  `json:"updateHandler" eh:"optional"`
    UpdatedHandler func (*Updated, *FileSystemService) (err error)  `json:"updatedHandler" eh:"optional"`
}

func (o *EventHandler) Apply(event eventhorizon.Event, entity eventhorizon.Entity) (err error) {
    switch event.EventType() {
    case CreateEvent:
        err = o.CreateHandler(event.Data().(*Create), entity.(*FileSystemService))
    case CreatedEvent:
        err = o.CreatedHandler(event.Data().(*Created), entity.(*FileSystemService))
    case DeleteEvent:
        err = o.DeleteHandler(event.Data().(*Delete), entity.(*FileSystemService))
    case DeletedEvent:
        err = o.DeletedHandler(event.Data().(*Deleted), entity.(*FileSystemService))
    case UpdateEvent:
        err = o.UpdateHandler(event.Data().(*Update), entity.(*FileSystemService))
    case UpdatedEvent:
        err = o.UpdatedHandler(event.Data().(*Updated), entity.(*FileSystemService))
    default:
		err = errors.New(fmt.Sprintf("Not supported event type '%v' for entity '%v", event.EventType(), entity))
	}
    return
}

func (o *EventHandler) SetupEventHandler() (err error) {

    //register event object factory
    eventhorizon.RegisterEventData(CreateEvent, func() eventhorizon.EventData {
		return &Create{}
	})

    //default handler implementation
    o.CreateHandler = func(event *Create, entity *FileSystemService) (err error) {
        //err = eh.EventHandlerNotImplemented(CreateEvent)
        return
    }

    //register event object factory
    eventhorizon.RegisterEventData(CreatedEvent, func() eventhorizon.EventData {
		return &Created{}
	})

    //default handler implementation
    o.CreatedHandler = func(event *Created, entity *FileSystemService) (err error) {
        if err = eh.ValidateNewId(entity.Id, event.Id, FileSystemServiceAggregateType); err == nil {
            entity.Name = event.Name
            entity.Id = event.Id
        }
        return
    }

    //register event object factory
    eventhorizon.RegisterEventData(DeleteEvent, func() eventhorizon.EventData {
		return &Delete{}
	})

    //default handler implementation
    o.DeleteHandler = func(event *Delete, entity *FileSystemService) (err error) {
        //err = eh.EventHandlerNotImplemented(DeleteEvent)
        return
    }

    //register event object factory
    eventhorizon.RegisterEventData(DeletedEvent, func() eventhorizon.EventData {
		return &Deleted{}
	})

    //default handler implementation
    o.DeletedHandler = func(event *Deleted, entity *FileSystemService) (err error) {
        if err = eh.ValidateIdsMatch(entity.Id, event.Id, FileSystemServiceAggregateType); err == nil {
            *entity = *NewFileSystemService()
        }
        return
    }

    //register event object factory
    eventhorizon.RegisterEventData(UpdateEvent, func() eventhorizon.EventData {
		return &Update{}
	})

    //default handler implementation
    o.UpdateHandler = func(event *Update, entity *FileSystemService) (err error) {
        //err = eh.EventHandlerNotImplemented(UpdateEvent)
        return
    }

    //register event object factory
    eventhorizon.RegisterEventData(UpdatedEvent, func() eventhorizon.EventData {
		return &Updated{}
	})

    //default handler implementation
    o.UpdatedHandler = func(event *Updated, entity *FileSystemService) (err error) {
        if err = eh.ValidateIdsMatch(entity.Id, event.Id, FileSystemServiceAggregateType); err == nil {
            entity.Name = event.Name
        }
        return
    }
    return
}


const FileSystemServiceAggregateType eventhorizon.AggregateType = "FileSystemService"

type AggregateInitializer struct {
    *eh.AggregateInitializer
    *CommandHandler
    *EventHandler
    ProjectorHandler *EventHandler `json:"projectorHandler" eh:"optional"`
}



func New@@EMPTY@@(eventStore eventhorizon.EventStore, eventBus eventhorizon.EventBus, eventPublisher eventhorizon.EventPublisher, 
                commandBus *bus.CommandHandler, 
                readRepos func (string, func () (ret eventhorizon.Entity) ) (ret eventhorizon.ReadWriteRepo) ) (ret *AggregateInitializer) {
    
    commandHandler := &FileSystemServiceCommandHandler{}
    eventHandler := &FileSystemServiceEventHandler{}
    entityFactory := func() eventhorizon.Entity { return NewFileSystemService() }
    ret = &AggregateInitializer{AggregateInitializer: eh.NewAggregateInitializer(FileSystemServiceAggregateType,
        func(id eventhorizon.UUID) eventhorizon.Aggregate {
            return eh.NewAggregateBase(FileSystemServiceAggregateType, id, commandHandler, eventHandler, entityFactory())
        }, entityFactory,
        FileSystemServiceCommandTypes().Literals(), FileSystemServiceEventTypes().Literals(), eventHandler,
        []func() error{commandHandler.SetupCommandHandler, eventHandler.SetupEventHandler},
        eventStore, eventBus, eventPublisher, commandBus, readRepos), FileSystemServiceCommandHandler: commandHandler, FileSystemServiceEventHandler: eventHandler, ProjectorHandler: eventHandler,
    }

    return
}


type FileSystemEventhorizonInitializer struct {
    eventStore eventhorizon.EventStore `json:"eventStore" eh:"optional"`
    eventBus eventhorizon.EventBus `json:"eventBus" eh:"optional"`
    eventPublisher eventhorizon.EventPublisher `json:"eventPublisher" eh:"optional"`
    commandBus *bus.CommandHandler `json:"commandBus" eh:"optional"`
    FileSystemServiceAggregateInitializer *AggregateInitializer `json:"fileSystemServiceAggregateInitializer" eh:"optional"`
}

func New@@EMPTY@@(eventStore eventhorizon.EventStore, eventBus eventhorizon.EventBus, eventPublisher eventhorizon.EventPublisher, 
                commandBus *bus.CommandHandler, 
                readRepos func (string, func () (ret eventhorizon.Entity) ) (ret eventhorizon.ReadWriteRepo) ) (ret *FileSystemEventhorizonInitializer) {
    fileSystemServiceAggregateInitializer := New@@EMPTY@@(eventStore, eventBus, eventPublisher, commandBus, readRepos)
    ret = &FileSystemEventhorizonInitializer{
        eventStore: eventStore,
        eventBus: eventBus,
        eventPublisher: eventPublisher,
        commandBus: commandBus,
        FileSystemServiceAggregateInitializer: fileSystemServiceAggregateInitializer,
    }
    return
}

func (o *FileSystemEventhorizonInitializer) Setup() (err error) {
    
    if err = o.FileSystemServiceAggregateInitializer.Setup(); err != nil {
        return
    }

    return
}









