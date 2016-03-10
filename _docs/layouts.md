Overview
---

termui offers two layout system: [Absolute]() and [Grid](). The two concept actually spawned from Web:

- The __Absolute layout__ is a plain coordination system, like [CSS position property](https://developer.mozilla.org/en/docs/Web/CSS/position) `position: absolute`. You will need manually assign `.X`, `.Y`, `.Width` and `.Height` to a component.
- The __Grid system__ actually is a simplified version of [the 12 columns CSS grid system](http://www.w3schools.com/bootstrap/bootstrap_grid_system.asp) on terminal. You do not need to bother setting positions and width properties, these values will be synced up according to their containers.

!!! note
	`Align` property can help you set your component position based on terminal window. Find more at [Magic Variables](#magic-variables)

__Cons and pros:__

- Use of Absolute layout gives you maximum control over how to arrange your components, while you have
to put a little more effort to set things up. Fortunately there are some "magic variables" may help you out.
- Grid layout can save you some time, it adjusts components location and size based on it's container. But note that you do need to set `.Height` property to each components because termui can not decide it for you.


Absolute Layout
---

Grid Layout
---

Magic Variables
---
