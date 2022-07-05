function trelay.on_connect(ctx)
	print("Wow, a connection!")
end

function trelay.on_player_packet(ctx)
	print("Got packet!")
	-- local pl = ctx.player

	-- if pl.Name() ~= "foobar" then
	-- 	pl:Disconnect("Your name is not 'foobar'!")
	-- end

	-- pl:ChangeServer("localhost:7878")
end
