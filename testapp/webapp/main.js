const getBoardRoot = () => {
  const root = document.querySelector('.checkers-board')
  if (root === null) {
    throw new Error('Cannot find board root: .checkers-board')
  }
  return root
}

const getBoardStatus = () => {
  const status = document.querySelector('.checkers-board + h2')
  if (status === null) {
    throw new Error('Cannot find board status: .checkers-board + h2')
  }
  return status
}

const getStartButton = () => {
  const button = document.querySelector('#btn-start-game')
  if (button === null) {
    throw new Error('Cannot find start button: #btn-start-game')
  }
  return button
}

const getBlackSquare = nth => {
  let index = 1
  for (const square of getBoardRoot().children) {
    if (square.classList.contains('red')) {
      continue
    }
    if (index === nth) {
      return square
    }
    index++
  }
  throw new Error(`Black square #${nth} was not found`)
}

const createBoard = () => {
  const root = getBoardRoot()
  for (let i = 0; i < 64; i++) {
    const square = document.createElement('div')
    const remainder = (~~(i / 8) + i) % 2
    square.classList.add(remainder === 0 ? 'red' : 'black')
    root.appendChild(square)
  }
}

const onSquareClick = handler => {
  let index = 1
  for (const square of getBoardRoot().children) {
    if (square.classList.contains('red')) {
      continue
    }
    let squareIndex = index
    square.onclick = event => {
      event.stopPropagation()
      handler(squareIndex)
    }
    index++
  }
}

const onBoardMove = handler => {
  let start = null
  const setStart = newStart => {
    if (start !== null) {
      getBlackSquare(start).classList.remove('chosen')
    }
    start = newStart
    if (start !== null) {
      getBlackSquare(start).classList.add('chosen')
    }
  }
  window.addEventListener('click', () => {
    setStart(null)
  })
  onSquareClick(square => {
    if (square === start) {
      setStart(null)
      return
    }
    if (start === null) {
      setStart(square)
      return
    }
    handler(start, square)
    setStart(null)
  })
}

const applyBoardState = ({ pieces, status }) => {
  getBoardStatus().textContent = `Status: ${status}`
  for (const { Color, Type, Square } of pieces) {
    const square = getBlackSquare(Square)
    for (const each of square.children) {
      each.remove()
    }
    const checker = document.createElement('div')
    checker.classList.add(
      'checker',
      Color ? 'red-checker' : 'black-checker',
      Type ? 'king-checker' : 'man-checker',
    )
    square.appendChild(checker)
  }
}

const startGame = async () => {
  const response = await fetch('api/start-game')
  const state = await response.json()
  applyBoardState(state)
}

const main = () => {
  createBoard()
  getStartButton().onclick = startGame
  onBoardMove((from, to) => {
    console.log(`from: ${from}, to: ${to}`)
  })
}

main()
