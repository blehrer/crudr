- use List bubble
- session state enum in main model like this:
    ```golang
    type sessionState int
    const (
       "/" = iota,
       "/paths",
       "/paths/{id}"
    )

    type MainModel struct{
       state sessionState 
       view1 View1Model
       view2 View2Model
       view3 View3Model
    }
    ```
- main model always gets messages first. it delegates downwards
