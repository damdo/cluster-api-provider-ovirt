package ovirtclient

import (
	"fmt"

	ovirtsdk "github.com/ovirt/go-ovirt"
)

func (o *oVirtClient) AddTagToVM(id string, tagID string, retries ...RetryStrategy) (err error) {
	retries = defaultRetries(retries, defaultWriteTimeouts())
	err = retry(
		fmt.Sprintf("adding tag %s to VM %s", tagID, id),
		o.logger,
		retries,
		func() error {
			_, err := o.conn.SystemService().VmsService().VmService(id).TagsService().Add().
				Tag(ovirtsdk.NewTagBuilder().Id(tagID).MustBuild()).Send()

			if err != nil {
				return err
			}
			return nil
		})
	return
}

func (m *mockClient) AddTagToVM(id string, tagID string, retries ...RetryStrategy) (err error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if _, ok := m.vms[id]; !ok {
		return newError(ENotFound, "VM with ID %s not found", id)
	}

	if _, ok := m.tags[tagID]; !ok {
		return newError(ENotFound, "tag with ID %s not found", tagID)
	}

	m.vms[id].tagIDs = append(m.vms[id].tagIDs, tagID)
	return nil

}

func (o *oVirtClient) AddTagToVMByName(id string, tagName string, retries ...RetryStrategy) (err error) {
	retries = defaultRetries(retries, defaultWriteTimeouts())
	err = retry(
		fmt.Sprintf("adding tag %s to VM %s", tagName, id),
		o.logger,
		retries,
		func() error {
			_, err := o.conn.SystemService().VmsService().VmService(id).TagsService().Add().
				Tag(ovirtsdk.NewTagBuilder().Name(tagName).MustBuild()).Send()

			return err
		})
	return
}

func (m *mockClient) AddTagToVMByName(id string, tagName string, retries ...RetryStrategy) (err error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if _, ok := m.vms[id]; !ok {
		return newError(ENotFound, "VM with ID %s not found", id)
	}

	for tagID, tag := range m.tags {
		if tag.name == tagName {
			m.vms[id].tagIDs = append(m.vms[id].tagIDs, tagID)
			return nil
		}
	}

	return newError(ENotFound, "Tag with Name %s not found", tagName)

}
