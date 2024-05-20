# Arena allocation

## Table of contents
{{ table of contents }}

## Code block index
{{ index }}

## What is an arena?

The arena is one of the simplest memory management strategies available. Here's what I know about it.

An **arena** (or region) is just a block of memory put aside at the start the program to do whatever with it. While every pointer allocated through the standard `malloc` assumes individual lifetimes, where each of them needs to be freed manually at some point, arena allocation can be taught of as a "dynamic stack". 

To begin with, here's the struct that holds the arena data. Not much, just a pointer to the buffer itself and variables to track its size and how much of it is used.

```c Arena struct
struct Arena
{
    uint8_t* area;
    size_t size;
    size_t used;
};
```

Allocating something is as simple as giving away the pointer to base + offset, where the offset is actually the number of bytes allocated. Not much else to keep track of. For a memory address let's say `0x8000`, when you allocate 5 bytes it will give you back the same memory address, but the next allocation will return `0x8005` as the address.

It's a "malloc"—but we are not setting aside some block of memory independently, rather we're "taking" free bytes away from the backing memory.

```c "malloc" function
void *take (struct Arena *a, size_t n)
{
    void *result;

    // See if we would've run out of
    // memory when allocating.
    if (a->used + n > a->size)
    {
        // If so, give up here.
        return NULL;
    }

    // Otherwise, increase the space used
    // and give a pointer to the beginning
    // of the space.
    result = a->area + a->used;
    a->used += n;

    return result;
}
```

"Free"ing is also possible, in some sense; in that the number of bytes taken is simply decreased. Because of the "stack"-like nature of these kinds of allocations, freeing a pointer allocated from the backing memory will not be valid unless *all* of the pointers succeeding it are freed, first.

```c "free" function
void give (struct Arena *a, size_t n)
{
    a->used -= n;
    return;
}
```

## A practical example

Let's imagine a program (like a game) that uses this system. The program manages a state and some amount of entities.

```c State struct
struct State
{
    struct Arena *arena;
    struct Entity *first_free;
};
```

```c Entity struct
struct Entity
{
    struct Entity *next;
    int x;
    int y;
};
```

```c Set memory base
uint8_t area[80] = { 0 };

/* simulate uninit'd memory */
for (size_t i = 0; i < sizeof(area); i++)
{
    area[i] = rand() % 0xff;
}
```

```c "calloc" function
void *takeZero (struct Arena *a, size_t n)
{
    void *res = take(a, n);
    if (res == NULL)
    {
        return res;
    }

    memset(res, 0, n);
    return res;
}
```

```c Entity allocation
struct Entity *alloc_entity(struct State *s)
{
    struct Entity *res = s->first_free;
    if (res == NULL)
    {
        /* initializes the Entity by zero-filling it */
        /* ->next == NULL */
        return takeZero(s->arena, sizeof(struct Entity));
    }
    s->first_free = s->first_free->next;
    return res;
}
```

```c Entity freeing
void release_entity(struct State *s, struct Entity *e)
{
    /* Technically the entity isn't freed */
    e->next = s->first_free;
    s->first_free = e;
}
```

```c The "free everything" function
void reset (struct Arena *a)
{
    a->used = 0;
}
```

```c Initialize arena and state
struct Arena a = {
    .area = area,
    .size = sizeof(area),
    .used = 0,
};

struct State s = {
    .arena = &a,
    .first_free = NULL,
};
```

```c Allocate and deallocate a couple of entities
struct Entity *entity_a = alloc_entity(&s);
struct Entity *entity_b = alloc_entity(&s);

release_entity(&s, entity_a);
release_entity(&s, entity_b);

struct Entity *entity_c = alloc_entity(&s);

release_entity(&s, entity_c);
```

### Final program

```c Entry point
int main (void)
{
    @{Set memory base}
    @{Initialize arena and state}
    @{Allocate and deallocate a couple of entities}
    return 0;
}
```

```c Includes
#include <stdint.h>
#include <stddef.h>
#include <string.h>
#include <stdlib.h>
```

```c /arenas.c
@{Includes}

/* Struct defines */
@{Arena struct}
@{Entity struct}
@{State struct}

/* Function defines */
@{"malloc" function}
@{"calloc" function}
@{"free" function}
@{Entity allocation}
@{Entity freeing}
@{The "free everything" function}

/* Main */
@{Entry point}
```

## References

* https://www.gingerbill.org/article/2019/02/15/memory-allocation-strategies-003/
* https://www.rfleury.com/p/untangling-lifetimes-the-arena-allocator
