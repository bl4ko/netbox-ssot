package inventory

func (nbi *NetBoxInventory) DeleteOrphans() error {
	nbi.Logger.Info("Deleting orphaned objects...")
	for objectAPIPath, ids := range nbi.OrphanManager {
		nbi.Logger.Info("Deleting orphaned objects of type ", objectAPIPath, " with IDs ", ids)
		err := nbi.NetboxApi.BulkDeleteObjects(objectAPIPath, ids)
		if err != nil {
			return err
		}
	}
	return nil
}
