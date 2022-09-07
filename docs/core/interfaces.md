# Interfaces
The core functionality is interfaced in `interfaces.go`, definitions of the interfaced functions can be found in the table below.

Type `IDer`:

| Function             | Definition                   |
|----------------------|------------------------------|
| GetID() *uuid.UUID*  | Gets the UUID for an entity  |

Type `Initializer`:

| Function                                  | Definition                  |
|-------------------------------------------|-----------------------------|
| Initialize(ctx context.Context) *error*   | Initializes a new instance  |

Type `Runner`:

| Function       | Definition       |
|----------------|------------------|
| Run() *error*  | Runs an instance |

Type `Stopper`:

| Function        | Definition        |
|-----------------|-------------------|
| Stop() *error*  | Stops an instance |

Type `Loader`:

| Function        | Definition        |
|-----------------|-------------------|
| Load() *error*  | Loads an instance |

Type `Saver`:

| Function       | Definition        |
|----------------|-------------------|
| Save() *error* | Saves an instance |

Type `APIRegister`:

| Function                          | Definition               |
|-----------------------------------|--------------------------|
| RegisterAPI(r *gin.Engine) *void* | Registers new API routes |

Type `Node`:

| Function                                    | Definition                                  |
|---------------------------------------------|---------------------------------------------|
| GetWorlds() *Worlds*                        | Gets all existing worlds for this node      |
| GetAssets2d() *Assets2d*                    | Gets all existing 2D assets for this node   |
| GetAssets3d() *Assets3d*                    | Gets all existing 3D assets for this node   |
| GetSpaceTypes() *SpaceTypes*                | Gets all existing Space Types for this node |
| AddAPIRegister(register APIRegister) *void* | <mark>?</mark>                              |

Type `Worlds`:

| Function                                            | Definition                                                                                                  |
|-----------------------------------------------------|-------------------------------------------------------------------------------------------------------------|
| GetWorld(worldID uuid.UUID) *(World, bool)*         | Gets an instance of a **World**, takes a **World** uuid, returns the **World** instance                     |
| GetWorlds(map[uuid.UUID]World) *void*               | Gets multiple instances of **World** based on a map of **World** uuids                                      |
| AddWorld(world World, updateDB bool) *error*        | Adds a **World** instance to the universe, takes a single **World** object, returns an optional error       |
| AddWorlds(worlds []World, updateDB bool) *error*    | Adds multiple **World** to the universe, takes a slice of **World** objects, returns an optional error      |
| RemoveWorld(world World, updateDB bool) *error*     | Removes a **World** from the universe, takes a single **World** object, returns an optional error           |
| RemoveWorlds(worlds []World, updateDB bool) *error* | Removes multiple **World** from the universe, takes a slice of **World** objects, returns an optional error |
 
Type `Space`:

| Function                                                                                | Definition                                                                                                                                                            |
|-----------------------------------------------------------------------------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| GetWorld() *World*                                                                      | Gets the **World** belonging to this **Space**, returns an updated **World**                                                                                          |
| GetParent() *Space*                                                                     | Gets the parent **Space** of this **Space**, returns the parent **Space**                                                                                             |
| SetParent(parent Space, updateDB bool) *error*                                          | Sets a new parent **Space** of this **Space**, takes a **Space** and an *updateDB* boolean to update the database record, returns an optional error                   |
| GetOwnerID() *uuid.UUID*                                                                | Gets the ownerID (User that has created the **Space**), returns a uuid                                                                                                |
| SetOwnerID(ownerID uuid.UUID, updateDB bool) *error*                                    | Sets the ownerID (User that has created the **Space**), takes a new ownedID (uuid) and an *updateDB* boolean to update the database record, returns an optional error |
| GetPosition() **cmath.Vec3*                                                             | Gets the position of the **Space**, returns a Vector3 position object                                                                                                 |
| SetPosition(position *cmath.Vec3, updateDB bool) *error*                                | Sets the position of the **Space**, takes a Vector3 position and an *updateDB* boolean to update the database record, returns an optional error                       |
| GetOptions() **entry.SpaceOptions*                                                      | Instance url                                                                                                                                                          |
| GetEffectiveOptions() **entry.SpaceOptions*                                             | Instance url                                                                                                                                                          |
| SetOptions(modifyFn modify.Fn[entry.SpaceOptions], updateDB bool) **entry.SpaceOptions* | Instance url                                                                                                                                                          |
| GetAsset2D() *Asset2d*                                                                  | Instance url                                                                                                                                                          |
| SetAsset2D(asset2d Asset2d, updateDB bool) *error*                                      | Instance url                                                                                                                                                          |
| GetAsset3D() *Asset3d*                                                                  | Instance url                                                                                                                                                          |
| SetAsset3D(asset3d Asset3d, updateDB bool) *error*                                      | Instance url                                                                                                                                                          |
| GetSpaceType() *SpaceType*                                                              | Instance url                                                                                                                                                          |
| SetSpaceType(spaceType SpaceType, updateDB bool) *error*                                | Instance url                                                                                                                                                          |
| GetEntry() **entry.Space*                                                               | Instance url                                                                                                                                                          |
| LoadFromEntry(entry *entry.Space, recursive bool) *error*                               | Instance url                                                                                                                                                          |
| GetSpace(spaceID uuid.UUID, recursive bool) *Space, bool*                               | Instance url                                                                                                                                                          |
| GetSpaces(recursive bool) *map[uuid.UUID]Space*                                         | Instance url                                                                                                                                                          |
| AddSpace(space Space, updateDB bool) *error*                                            | Instance url                                                                                                                                                          |
| AddSpaces(spaces []Space, updateDB bool) *error*                                        | Instance url                                                                                                                                                          |
| RemoveSpace(space Space, recursive, updateDB bool) *bool, error*                        | Instance url                                                                                                                                                          |
| RemoveSpaces(spaces []Space, recursive, updateDB bool) *bool, error*                    | Instance url                                                                                                                                                          |
| GetUser(userID uuid.UUID, recursive bool) *User, bool*                                  | Instance url                                                                                                                                                          |
| GetUsers(recursive bool) *map[uuid.UUID]User*                                           | Instance url                                                                                                                                                          |
| AddUser(user User, updateDB bool) *error*                                               | Instance url                                                                                                                                                          |
| RemoveUser(user User, updateDB bool) *error*                                            | Instance url                                                                                                                                                          |
| SendToUser(userID uuid.UUID, msg *websocket.PreparedMessage, recursive bool) *error*    | Instance url                                                                                                                                                          |
| Broadcast(msg *websocket.PreparedMessage, recursive bool) *error*                       | Instance url                                                                                                                                                          |

Type `User`:

| Function                                     | Definition   |
|----------------------------------------------|--------------|
| GetWorld() *World*                           | Instance url |
| SetWorld(world World, updateDB bool) *error* | Instance url |
| GetSpace() *Space*                           | Instance url |
| SetSpace(space Space, updateDB bool) *error* | Instance url |

Type `SpaceTypes`:

| Function                                                         | Definition   |
|------------------------------------------------------------------|--------------|
| GetSpaceType(spaceTypeID uuid.UUID) *SpaceType, bool*            | Instance url |
| GetSpaceTypes() *map[uuid.UUID]SpaceType*                        | Instance url |
| AddSpaceType(spaceType SpaceType, updateDB bool) *error*         | Instance url |
| AddSpaceTypes(spaceTypes []SpaceType, updateDB bool) *error*     | Instance url |
| RemoveSpaceType(spaceType SpaceType, updateDB bool) *error*      | Instance url |
| RemoveSpaceTypes(spaceTypes []SpaceType, updateDB bool) *error*  | Instance url |

Type `SpaceType`:

| Function                                                                  | Definition   |
|---------------------------------------------------------------------------|--------------|
| GetName() *string*                                                        | Instance url |
| SetName(name string, updateDB bool) *error*                               | Instance url |
| GetCategoryName() *string*                                                | Instance url |
| SetCategoryName(categoryName string, updateDB bool) *error*               | Instance url |
| GetDescription() **string*                                                | Instance url |
| SetDescription(description *string, updateDB bool) *error*                | Instance url |
| GetOptions() **entry.SpaceOptions*                                        | Instance url |
| SetOptions(modifyFn modify.Fn[entry.SpaceOptions], updateDB bool) *error* | Instance url |
| GetAsset2d() *Asset2d*                                                    | Instance url |
| SetAsset2d(asset2d Asset2d, updateDB bool) *error*                        | Instance url |
| GetAsset3d() *Asset3d*                                                    | Instance url |
| SetAsset3d(asset3d Asset3d, updateDB bool) *error*                        | Instance url |
| GetEntry() **entry.SpaceType*                                             | Instance url |
| LoadFromEntry(entry *entry.SpaceType) *error*                             | Instance url |

Type `Assets2d`:

| Function                                             | Definition   |
|------------------------------------------------------|--------------|
| GetAsset2d(asset2dID uuid.UUID) **                   | Instance url |
| GetAssets2d() **                                     | Instance url |
| AddAsset2d(asset2d Asset2d, updateDB bool) **        | Instance url |
| AddAssets2d(assets2d []Asset2d, updateDB bool) **    | Instance url |
| RemoveAsset2d(asset2d Asset2d, updateDB bool) **     | Instance url |
| RemoveAssets2d(assets2d []Asset2d, updateDB bool) ** | Instance url |

Type `Asset2d`:

| Function                                                                    | Definition   |
|-----------------------------------------------------------------------------|--------------|
| GetName() *string*                                                          | Instance url |
| SetName(name string, updateDB bool) *error*                                 | Instance url |
| GetOptions() **entry.Asset2dOptions*                                        | Instance url |
| SetOptions(modifyFn modify.Fn[entry.Asset2dOptions], updateDB bool) *error* | Instance url |
| GetEntry() **entry.Asset2d*                                                 | Instance url |
| LoadFromEntry(entry *entry.Asset2d) *error*                                 | Instance url |

Type `Assets3d`:

| Function                                                  | Definition   |
|-----------------------------------------------------------|--------------|
| GetAsset3d(asset3dID uuid.UUID) *Asset3d, bool*           | Instance url |
| GetAssets3d() *map[uuid.UUID]Asset3d*                     | Instance url |
| AddAsset3d(asset3d Asset3d, updateDB bool) *error*        | Instance url |
| AddAssets3d(assets3d []Asset3d, updateDB bool) *error*    | Instance url |
| RemoveAsset3d(asset3d Asset3d, updateDB bool) *error*     | Instance url |
| RemoveAssets3d(assets3d []Asset3d, updateDB bool) *error* | Instance url |

Type `Asset3d`:

| Function                                                                    | Definition   |
|-----------------------------------------------------------------------------|--------------|
| GetName() *string*                                                          | Instance url |
| SetName(name string, updateDB bool) *error*                                 | Instance url |
| GetOptions() **entry.Asset3dOptions*                                        | Instance url |
| SetOptions(modifyFn modify.Fn[entry.Asset3dOptions], updateDB bool) *error* | Instance url |
| GetEntry() **entry.Asset3d*                                                 | Instance url |
| LoadFromEntry(entry *entry.Asset3d) *error*                                 | Instance url |