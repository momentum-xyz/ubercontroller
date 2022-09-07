# Interfaces
The core functionality is interfaced in `interfaces.go`, definitions of the interfaced functions can be found in the table below.

Type `IDer`:

| Function             | Definition                   |
|----------------------|------------------------------|
| GetID() *uuid.UUID*  | Gets the UUID for an entity  |

Type `Initializer`:

| Function                                  | Definition   |
|-------------------------------------------|--------------|
| Initialize(ctx context.Context) *error*   | Instance url |

Type `Runner`:

| Function       | Definition   |
|----------------|--------------|
| Run() *error*  | Instance url |

Type `Stopper`:

| Function        | Definition   |
|-----------------|--------------|
| Stop() *error*  | Instance url |

Type `Loader`:

| Function        | Definition   |
|-----------------|--------------|
| Load() *error*  | Instance url |

Type `Saver`:

| Function       | Definition   |
|----------------|--------------|
| Save() *error* | Instance url |

Type `APIRegister`:

| Function                          | Definition   |
|-----------------------------------|--------------|
| RegisterAPI(r *gin.Engine) *void* | Instance url |

Type `Node`:

| Function                                    | Definition   |
|---------------------------------------------|--------------|
| GetWorlds() *Worlds*                        | Instance url |
| GetAssets2d() *Assets2d*                    | Instance url |
| GetAssets3d() *Assets3d*                    | Instance url |
| GetSpaceTypes() *SpaceTypes*                | Instance url |
| AddAPIRegister(register APIRegister) *void* | Instance url |

Type `Worlds`:

| Function                                            | Definition   |
|-----------------------------------------------------|--------------|
| GetWorld(worldID uuid.UUID) *(World, bool)*         | Instance url |
| GetWorlds(map[uuid.UUID]World) *void*               | Instance url |
| AddWorld(world World, updateDB bool) *error*        | Instance url |
| AddWorlds(worlds []World, updateDB bool) *error*    | Instance url |
| RemoveWorld(world World, updateDB bool) *error*     | Instance url |
| RemoveWorlds(worlds []World, updateDB bool) *error* | Instance url |
 
Type `Space`:

| Function                                                                                | Definition   |
|-----------------------------------------------------------------------------------------|--------------|
| GetWorld() *World*                                                                      | Instance url |
| GetParent() *Space*                                                                     | Instance url |
| SetParent(parent Space, updateDB bool) *error*                                          | Instance url |
| GetOwnerID() *uuid.UUID*                                                                | Instance url |
| SetOwnerID(ownerID uuid.UUID, updateDB bool) *error*                                    | Instance url |
| GetPosition() **cmath.Vec3*                                                             | Instance url |
| SetPosition(position *cmath.Vec3, updateDB bool) *error*                                | Instance url |
| GetOptions() **entry.SpaceOptions*                                                      | Instance url |
| GetEffectiveOptions() **entry.SpaceOptions*                                             | Instance url |
| SetOptions(modifyFn modify.Fn[entry.SpaceOptions], updateDB bool) **entry.SpaceOptions* | Instance url |
| GetAsset2D() *Asset2d*                                                                  | Instance url |
| SetAsset2D(asset2d Asset2d, updateDB bool) *error*                                      | Instance url |
| GetAsset3D() *Asset3d*                                                                  | Instance url |
| SetAsset3D(asset3d Asset3d, updateDB bool) *error*                                      | Instance url |
| GetSpaceType() *SpaceType*                                                              | Instance url |
| SetSpaceType(spaceType SpaceType, updateDB bool) *error*                                | Instance url |
| GetEntry() **entry.Space*                                                               | Instance url |
| LoadFromEntry(entry *entry.Space, recursive bool) *error*                               | Instance url |
| GetSpace(spaceID uuid.UUID, recursive bool) *Space, bool*                               | Instance url |
| GetSpaces(recursive bool) *map[uuid.UUID]Space*                                         | Instance url |
| AddSpace(space Space, updateDB bool) *error*                                            | Instance url |
| AddSpaces(spaces []Space, updateDB bool) *error*                                        | Instance url |
| RemoveSpace(space Space, recursive, updateDB bool) *bool, error*                        | Instance url |
| RemoveSpaces(spaces []Space, recursive, updateDB bool) *bool, error*                    | Instance url |
| GetUser(userID uuid.UUID, recursive bool) *User, bool*                                  | Instance url |
| GetUsers(recursive bool) *map[uuid.UUID]User*                                           | Instance url |
| AddUser(user User, updateDB bool) *error*                                               | Instance url |
| RemoveUser(user User, updateDB bool) *error*                                            | Instance url |
| SendToUser(userID uuid.UUID, msg *websocket.PreparedMessage, recursive bool) *error*    | Instance url |
| Broadcast(msg *websocket.PreparedMessage, recursive bool) *error*                       | Instance url |

Type `User`:

| Function | Definition   |
|----------|--------------|
| url      | Instance url |

Type `SpaceTypes`:

| Function | Definition   |
|----------|--------------|
| url      | Instance url |

Type `SpaceType`:

| Function | Definition   |
|----------|--------------|
| url      | Instance url |

Type `Assets2d`:

| Function | Definition   |
|----------|--------------|
| url      | Instance url |

Type `Asset2d`:

| Function | Definition   |
|----------|--------------|
| url      | Instance url |

Type `Assets3d`:

| Function | Definition   |
|----------|--------------|
| url      | Instance url |

Type `Asset3d`:

| Function | Definition   |
|----------|--------------|
| url      | Instance url |