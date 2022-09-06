# Interfaces
The core functionality is interfaced in `interfaces.go`, definitions of the interfaced functions can be found in the table below.

Type `IDer`:

| Function | Definition                   |
|----------|------------------------------|
| GetID()  | Gets the UUID for an entity  |

Type `Initializer`:

| Function     | Definition   |
|--------------|--------------|
| Initialize() | Instance url |

Type `Runner`:

| Function | Definition   |
|----------|--------------|
| Run()    | Instance url |

Type `Stopper`:

| Function | Definition   |
|----------|--------------|
| Stop()   | Instance url |

Type `Loader`:

| Function | Definition   |
|----------|--------------|
| Load()   | Instance url |

Type `Saver`:

| Function | Definition   |
|----------|--------------|
| Save()   | Instance url |

Type `APIRegister`:

| Function                    | Definition   |
|-----------------------------|--------------|
| RegisterAPI(r *gin.Engine)  | Instance url |

Type `Node`:

| Function                             | Definition   |
|--------------------------------------|--------------|
| GetWorlds()                          | Instance url |
| GetAssets2d()                        | Instance url |
| GetAssets3d()                        | Instance url |
| GetSpaceTypes()                      | Instance url |
| AddAPIRegister(register APIRegister) | Instance url |

Type `Worlds`:

| Function    | Definition   |
|-------------|--------------|
| GetWorld()  | Instance url |
| GetWorlds() | Instance url |
| GetWorld()  | Instance url |
| GetWorld()  | Instance url |
| GetWorld()  | Instance url |
| GetWorld()  | Instance url |

Type `Space`:

| Function | Definition   |
|----------|--------------|
| url      | Instance url |

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