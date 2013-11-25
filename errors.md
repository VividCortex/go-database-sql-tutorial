---
layout: page
permalink: /errors/
title: Error Handling
tags: 
image:
  feature: abstract-5.jpg
share: true
---

Go defines a special error constant, called `sql.ErrNoRows`, which is returned from `QueryRow()` when the
result is empty. This needs to be handled as a special case in most circumstances. An empty result is often
not considered an error by application code, and if you don't check whether an error is this special constant,
you'll cause application-code errors you didn't expect.

One might ask why an empty result set is considered an error. There's nothing erroneous about an empty set,
after all. The reason is that the `QueryRow()` method needs to use this special-case in order to let the caller
distinguish whether `QueryRow()` in fact found a row; without it, `Scan()` wouldn't do anything and you might
not realize that your variable didn't get any value from the database after all.


{% include toc.md %}
