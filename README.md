# GynDNS

GynDNS is an open source, self-hosted solution for dynamic DNS.
It is compatible with noip and dyndns clients, as long as you manage to point them to wherever you host gyndns.

# Project status

This is currently a **WIP** project. I don't recommend using it (yet).

## How it works

You set up a list of users (and passwords), and a list of hostnames they are able to update.
Then, the server listens for authenticated HTTP requests from these users, and save in memory
the address they want to associate with a given hostname.

A dns server is set up to respond to `A` queries for these hostnames.

## Setup

GynDNS can work as an standalone server, but I recommend setting it up behind a strong reverse
proxy that supports HTTPs, and using the DNS server as an slave for bind.
