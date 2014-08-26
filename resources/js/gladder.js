
Ranker = (function () {
	var r;
	r = {
		numeric_from_ladder: function (lrank) {
			if (lrank > 225 || lrank < 0) {
				return undefined;
			}
			if (lrank < 81) {
				return (lrank / 10 - 8);
			} else if (lrank < 155) {
				return ((lrank - 80) / 5);
			}
			return ((lrank - 155) / 2 + 15);
		},
		display_from_numeric: function (n) {
			if (n === undefined) {
				return undefined;
			}
			if (n <= 0) {
				return (-n + 1).toFixed(1) + 'd';
			}
			return (n).toFixed(1) + 'k';
		},
		display_from_ladder: function (lrank) {
			return r.display_from_numeric(r.numeric_from_ladder(lrank));
		}
	};
	r.ladder_from_display = function (disp_rank) {
		var match = /^\s*([0-9.]+)\s*([kKdD])\s*/.exec(disp_rank);
		if (match === null) {
			return null;
		}
		var val = parseFloat(match[1]);
		if (match[2].toLowerCase() === 'd') {
			return Math.round((80 - (val - 1) * 10));
		} else {
			if (val < 15) {
				return Math.round(80 + 5 * val);
			}
			return Math.round(155 + 2 * (val - 15));
		}
		return 0;
	}
	for (var i = 0; i <= 225; i++) {
		var d = r.display_from_ladder(i);
		var d2 = r.ladder_from_display(' ' + d + ' ');
		if (i !== d2) {
			throw(["failed test", i, d2]);
		}
	}

	return r
})();

function show_handicap(user_x, user_y, el) {
	// 1 is higher, and therefore white
	var user1, user2;
	if (user_x.rank < user_y.rank) {
		user1 = user_x;
		user2 = user_y;
	} else {
		user1 = user_y;
		user2 = user_x;
	}

	var calc_handicap = function (bs, rdiff) {
		var calc = function (rdiff, opts) {
			var passes = rdiff / opts.rdiff_per_stone;
			var complete_passes = Math.floor(passes + 0.001);
			var stones = complete_passes + 1;
			if (complete_passes == 0) {
				stones = 0;
			}
			if (stones > opts.stone_cap) {
				stones = opts.stone_cap; // handle rest with komi
			}
			var komi_passes = passes - complete_passes;
			var raw_komi = opts.base_komi * (1 - 2*komi_passes);
			var komi = Math.floor(raw_komi) + 0.5;
			console.log(opts);
			console.log({rdiff:rdiff, passes:passes, complete_passes:complete_passes, stones:stones, komi_passes:komi_passes, raw_komi: raw_komi, komi:komi});
			return {komi: komi, stones: stones};
		}
		var opts;
		if (bs === '9x9') {
			opts = {
				base_komi: 7.5,
				stone_cap: 6,
				rdiff_per_stone: 6.0,
			};
		} else if (bs === '13x13') {
			opts = {
				base_komi: 6.5,
				stone_cap: 9,
				rdiff_per_stone: 3.0,
			};
		} else if (bs === '19x19') {
			opts = {
				base_komi: 6.5,
				stone_cap: 12,
				rdiff_per_stone: 1.0,
			};
		}
		return calc(rdiff, opts);
	};
	var per_board = [];
	var rank_diff = Ranker.numeric_from_ladder(user2.rank) - Ranker.numeric_from_ladder(user1.rank);
	var board_config = [{ size: '9x9' }, {size: '13x13' }, {size: '19x19' }];
	$.each(board_config, function () {
		var bc = this;

		var hcap = calc_handicap(bc.size, rank_diff);
		var handicap = {
			board_size: bc.size,
			stones_to_black: hcap.stones,
			white: user1,
			white_komi: hcap.komi > 0 ? hcap.komi : 0,
			black: user2,
			black_komi: hcap.komi < 0 ? -hcap.komi: 0
		};
		per_board.push(handicap);
	});

	var dat = {
		per_board: per_board,
		black_gets: function () {
			var parts = [];
			if (this.stones_to_black > 0.0001) {
				parts.push(this.stones_to_black.toFixed() + " stones");
			}
			if (this.black_komi > 0.0001) {
				parts.push(this.black_komi.toFixed(1) + " points");
			}
			if (parts.length === 0) {
				return '';
			}
			return 'and gets ' + parts.join(' and ');
		},
		white_gets: function () {
			if (this.white_komi > 0.0001) {
				return 'and gets ' + this.white_komi.toFixed(1) + " points";
			}
			return '';
		}
	}

	var html = Mustache.render('<button class="reset">reset</button>' +
							   '<table class="handicap">{{#per_board}}<tr><td>On {{board_size}}</td>' +
							   '<td>' +
							   '<span class="black_player"><span class="plays">{{black.name}} plays black</span> {{black_gets}}</span><br>' +
							   '<span class="white_player"><span class="plays">{{white.name}} plays white</span> {{white_gets}}</span>' +
							   '</td>' +
							   '</tr>{{/per_board}}</table>' +
							   '', dat);
	el.html(html);
}

$(function () {

	$("table.ladder tr").each(function (e) {
		var lrank = parseInt($(this).data('rank'), 10);
		$(this).find("td.rank").html(Ranker.display_from_ladder(lrank));
	});
	
	var selected_users = [];
	var reset_selection = function() {
			$("table.ladder tr").removeClass("selected");
			$("#info .handicap").hide();
			$("#info .help").show();
			selected_users = [];
	}

	$("table.ladder").on('click', "td.name", function (e) {
		var el = $(e.target).closest('tr');
		if (selected_users.length > 1) {
			reset_selection();
		}
		el.addClass("selected");
		selected_users.push({
			name: el.data("name"),
			rank: el.data("rank")
		});
		if (selected_users.length === 2) {
			show_handicap(selected_users[0], selected_users[1], $("#info .handicap").show());
			$("#info .help").hide();
		}
	});

	$("table.ladder").on('click', "td.rank", function (e) {
		var tr = $(e.target).closest('tr');
		var lrank = tr.data('rank');
		var disp = Ranker.display_from_ladder(lrank);
		var new_rank = window.prompt("Alter rank for player " + tr.data('name'), disp);
		if (new_rank === null) {
			return;
		}
		var new_lrank = Ranker.ladder_from_display(new_rank);
		var url = '/player/' + tr.data('name') + '/';
		$.ajax({
			url: url,
			type: 'POST',
			data: { rank: new_lrank },
			success: function () {
				document.location.reload(true);
			},
			error: function (xhr, status) {
			}
		});
	});

	$(".handicap").on('click', '.reset', function () { reset_selection(); });
});
