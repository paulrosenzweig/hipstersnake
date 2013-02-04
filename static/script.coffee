ws = board = myName = theirName = undefined

connect = ->
  ws = new WebSocket("ws://#{document.location.host}/ws")
  ws.onopen = -> $("#connection-status").text("Connected! Waiting for another player...")
  ws.onerror = (e) ->
    console.log("error", e)
    ws = undefined
  ws.onclose = (e) ->
    console.log("close", e)
    ws = undefined
  ws.onmessage = (d) ->
    data = JSON.parse(d.data)
    if data.ping
      send(ping: "ping")
    else if data.myName
      myName = data.myName
      theirName = data.theirName
    else if data.hasEnded
      message =
        if data.playerOne.hasLost and data.playerTwo.hasLost
          "You both lost at the same time."
        else if data[myName].hasLost
          "You lost!"
        else
          "You won!"
      $("#connection-status").text(message).slideDown()
      $("#replay").fadeIn()
      ws.close()
    else if data.playerOne
      board ||= new Board(data)
      $("#connection-status").slideUp()
      $("#grid").fadeIn()
      board.update(data)

connect()

$("#replay").click replay = ->
  return unless ws is undefined
  $("#connection-status").text("Connecting...")
  board.clear()
  connect()
  $("#replay").fadeOut()

@Board = class Board
  constructor: (data) ->
    @width = data.width
    @height = data.height
    grid = $("#grid")
    @domElements = {}
    for y in [0...@height]
      for x in [0...@width]
        elem = $("<div>", class: "cell")
        elem.data({x, y})
        @domElements[[x,y]] = elem
        grid.append(elem)
      grid.append($("<div>", style: "clear: both;"))
    grid.on('touchstart', (evt) -> evt.preventDefault())
    grid.on('touchend', @clickHandler)
    grid.click(@clickHandler)
    @formerTails = []

  update: (data) ->
    me = data[myName]
    them = data[theirName]
    myTail = me.position[me.position.length-1]
    theirTail = them.position[them.position.length-1]
    if @formerTails.length is 0
      for pos in me.position
        @domElements[pos].toggleClass("me", true)
      for pos in them.position
        @domElements[pos].toggleClass("them", true)
    else
      for pos in @formerTails
        if pos.toString() not in [myTail.toString(), theirTail.toString()]
          @domElements[pos].removeClass("me them")
      @domElements[me.position[0]].removeClass("food").addClass("me")
      @domElements[them.position[0]].removeClass("food").addClass("them")
      for pos in data.food
        @domElements[pos].addClass("food")
    [@head, @tail...] = me.position
    @formerTails = [myTail, theirTail]

  clear: ->
    for pos, elem of @domElements
      elem.removeClass("food me them")

  heading: ->
    [@head[0] - @tail[0][0], @head[1] - @tail[0][1]]

  clickHandler: (evt) =>
    evt.preventDefault()
    {x, y} = $(evt.target).data()
    heading = @heading()
    dx = x - @head[0]
    dy = y - @head[1]
    
    moreX =  Math.abs(dx) > Math.abs(dy)

    # going from horizontal to vertical
    if heading[0] and dy
      send heading: if dy > 0 then 'down' else 'up'
      # about face
      if moreX and heading[0] * dx < 0
        send heading: if dx > 0 then 'right' else 'left'

    # going from vertical to horizontal
    else if heading[1] and dx  # moving vertically
      send heading: if dx > 0 then 'right' else 'left'
      # about face
      if not moreX and heading[1] * dy < 0   # clicked behind
        send heading: if dy > 0 then 'down' else 'up'

    return false

send = (m) ->
  console.log(m) unless m.ping
  ws.send(JSON.stringify(m)) unless ws is undefined

$("body").keydown (e) ->
  heading = switch e.keyCode
    when 37, 72 then 'left'
    when 38, 75 then 'up'
    when 39, 76 then 'right'
    when 40, 74 then 'down'
  send(heading: heading) unless heading is undefined

  if e.keyCode is 13 or e.keyCode is 32
    replay()
