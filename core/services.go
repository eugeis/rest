package core

import "os"

func (o *Eye) reloadServiceFactory() {
	if len(o.config.ExportFolder) > 0 {
		os.MkdirAll(o.config.ExportFolder, 0777)
	}

	oldServiceFactory := o.serviceFactory
	o.serviceFactory = o.buildServiceFactory()
	if oldServiceFactory != nil {
		oldServiceFactory.Close()
	}

	o.checks = make(map[string]Check)

	//register queries
	o.registerMultiPing()
	o.registerValidateChecks()
	o.registerMultiValidates()
	o.registerCompares()

	//register exporters
	o.exporters = make(map[string]Exporter)
	o.registerExporters()

}

func (o *Eye) buildServiceFactory() Factory {
	serviceFactory := NewFactory()
	for _, item := range o.config.MySql {
		serviceFactory.Add(&MySqlService{mysql: item, accessFinder: o.accessFinder})
	}

	for _, item := range o.config.Http {
		serviceFactory.Add(&HttpService{http: item, accessFinder: o.accessFinder})
	}

	for _, item := range o.config.Fs {
		serviceFactory.Add(&FsService{Fs: item})
	}

	for _, item := range o.config.Ps {
		serviceFactory.Add(&PsService{Ps: item})
	}

	for _, item := range o.config.Elastic {
		serviceFactory.Add(&ElasticService{elastic: item})
	}
	return serviceFactory
}
